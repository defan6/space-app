package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/defan6/space-app/inventory-service/pkg/models"
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort        = 50052
	httpPort        = 8083
	shutdownTimeout = 5 * time.Second
)

type Cache struct {
	storage map[uuid.UUID]*models.Part
	mu      *sync.RWMutex
}

type FilterChain struct {
	filters []Filter
}

func NewFilterChain(filters ...Filter) *FilterChain {
	fs := make([]Filter, 0, 5)

	fs = append(fs, filters...)

	return &FilterChain{
		filters: fs,
	}
}

func (fc *FilterChain) DoPartFilter(parts []*models.Part) []*models.Part {
	res := make([]*models.Part, 0, len(parts))
	for _, f := range fc.filters {
		res = f(parts)
	}
	return res
}

type Filter func([]*models.Part) []*models.Part

func UUIDFilter(ids []string) Filter {
	idsMap := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idsMap[id] = struct{}{}
	}
	return func(parts []*models.Part) []*models.Part {
		for _, t := range ids {
			idsMap[t] = struct{}{}
		}
		if ids == nil {
			return parts
		}
		res := make([]*models.Part, 0, len(parts))
		for _, p := range parts {
			if _, ok := idsMap[p.ID.String()]; ok {
				res = append(res, p)
			}
		}

		return res
	}
}

func CategoryFilter(categories []inventoryV1.Category) Filter {
	categoriesMap := make(map[models.Category]struct{}, len(categories))
	for _, cat := range categories {
		categoriesMap[CategoryFromProtoToModel(cat)] = struct{}{}
	}
	return func(parts []*models.Part) []*models.Part {
		if categories == nil {
			return parts
		}
		res := make([]*models.Part, 0, len(parts))

		for _, p := range parts {
			if _, ok := categoriesMap[p.Category]; ok {
				res = append(res, p)
			}
		}

		return res
	}
}

func NameFilter(names []string) Filter {
	namesMap := make(map[string]struct{}, len(names))
	for _, n := range names {
		namesMap[n] = struct{}{}
	}

	return func(parts []*models.Part) []*models.Part {
		if names == nil {
			return parts
		}
		res := make([]*models.Part, 0, len(parts))

		for _, p := range parts {
			if _, ok := namesMap[p.Name]; ok {
				res = append(res, p)
			}
		}

		return res
	}
}

func TagCategoryFilter(tags []string) Filter {
	tagsMap := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagsMap[t] = struct{}{}
	}
	return func(parts []*models.Part) []*models.Part {
		if tags == nil {
			return parts
		}
		res := make([]*models.Part, 0, len(parts))
		for _, p := range parts {
			for _, pt := range p.Tags {
				if _, ok := tagsMap[pt]; ok {
					res = append(res, p)
				}
			}
		}

		return res
	}
}

func ManufacturerCountryFilter(countries []string) Filter {
	countriesMap := make(map[string]struct{}, len(countries))
	for _, c := range countries {
		countriesMap[c] = struct{}{}
	}
	return func(parts []*models.Part) []*models.Part {
		if countries == nil {
			return parts
		}
		res := make([]*models.Part, 0, len(parts))
		for _, p := range parts {
			if _, ok := countriesMap[p.Manufacturer.Country]; ok {
				res = append(res, p)
			}
		}

		return res
	}
}

func NewCache() *Cache {
	return &Cache{
		storage: make(map[uuid.UUID]*models.Part),
		mu:      &sync.RWMutex{},
	}
}

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage *Cache
}

func NewInventoryService(storage *Cache) *InventoryService {
	return &InventoryService{
		storage: storage,
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("error listen port: %d", grpcPort)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listen server on port: %d", grpcPort)
		}
	}()
	s := grpc.NewServer()
	storage := NewCache()
	service := NewInventoryService(storage)
	reflection.Register(s)
	inventoryV1.RegisterInventoryServiceServer(s, service)

	go func() {
		log.Printf("grpc server started on port: %d", grpcPort)
		if err = s.Serve(lis); err != nil {
			log.Printf("failed to serve on port %d", grpcPort)
		}
	}()

	var gatewayServer *http.Server

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register grpc gateway: %v\n", err)
			return
		}

		fileServer := http.FileServer(http.Dir("../shared/api/inventory/v1/swagger"))

		httpMux := http.NewServeMux()

		httpMux.Handle("/api/v1/inventory", mux)

		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/inventory.swagger.json", fileServer)

		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, r)
		}))

		gatewayServer = &http.Server{
			Addr:              fmt.Sprintf("localhost:%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}
		log.Printf("http server with grpc gateway listening on port %d\n", httpPort)
		err = gatewayServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to serve http: %d\n", httpPort)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Stopping servers...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	done := make(chan struct{})

	go func() {
		if gatewayServer != nil {
			err = gatewayServer.Shutdown(ctx)
			if err != nil {
				log.Printf("failed to shutdown http server on port %d, %v\n", httpPort, err)
			}
		}
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("planning stopped grpc server on port: %d", grpcPort)
	case <-ctx.Done():
		log.Printf("Server cannot stopped. Terminating")
	}
}

func (s *InventoryService) GetPart(_ context.Context, r *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	if err := r.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()
	stringUUID := r.Uuid
	UUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %s", stringUUID)
	}
	part, ok := s.storage.storage[UUID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", UUID.String())
	}
	return &inventoryV1.GetPartResponse{
		Uuid:          part.ID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryFromModelToProto(part.Category),
		Dimensions:    DimensionsFromModelToProto(part.Dimensions),
		Manufacturer:  ManufacturerFromModelToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataFromModelToProto(part.Metadata),
		CreatedAt:     TimePtrToProtoTimestamp(part.CreatedAt),
		UpdatedAt:     TimePtrToProtoTimestamp(part.UpdatedAt),
	}, nil
}

func (s *InventoryService) GetParts(_ context.Context, r *inventoryV1.GetPartsRequest) (*inventoryV1.GetPartsResponse, error) {
	if err := r.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	parts := make([]*models.Part, 0, 10)

	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	for _, p := range s.storage.storage {
		parts = append(parts, p)
	}
	protoParts := make([]*inventoryV1.GetPartResponse, 0, len(parts))
	response := &inventoryV1.GetPartsResponse{}
	if r.PartsFilter == nil {
		for _, p := range parts {
			protoPart := &inventoryV1.GetPartResponse{
				Uuid:          p.ID.String(),
				Name:          p.Name,
				Description:   p.Description,
				Price:         p.Price,
				StockQuantity: p.StockQuantity,
				Category:      CategoryFromModelToProto(p.Category),
				Dimensions:    DimensionsFromModelToProto(p.Dimensions),
				Manufacturer:  ManufacturerFromModelToProto(p.Manufacturer),
				Tags:          p.Tags,
				Metadata:      MetadataFromModelToProto(p.Metadata),
				CreatedAt:     TimePtrToProtoTimestamp(p.CreatedAt),
				UpdatedAt:     TimePtrToProtoTimestamp(p.UpdatedAt),
			}
			protoParts = append(protoParts, protoPart)
		}
		response.GetPartsResponse = protoParts
		return response, nil
	}

	filterChain := NewFilterChain(
		UUIDFilter(r.PartsFilter.Uuids),
		CategoryFilter(r.PartsFilter.Categories),
		TagCategoryFilter(r.PartsFilter.Tags),
		ManufacturerCountryFilter(r.PartsFilter.ManufacturerCountries),
		NameFilter(r.PartsFilter.Names),
	)
	filteredParts := filterChain.DoPartFilter(parts)
	for _, fp := range filteredParts {
		protoPart := &inventoryV1.GetPartResponse{
			Uuid:          fp.ID.String(),
			Name:          fp.Name,
			Description:   fp.Description,
			Price:         fp.Price,
			StockQuantity: fp.StockQuantity,
			Category:      CategoryFromModelToProto(fp.Category),
			Dimensions:    DimensionsFromModelToProto(fp.Dimensions),
			Manufacturer:  ManufacturerFromModelToProto(fp.Manufacturer),
			Tags:          fp.Tags,
			Metadata:      MetadataFromModelToProto(fp.Metadata),
			CreatedAt:     TimePtrToProtoTimestamp(fp.CreatedAt),
			UpdatedAt:     TimePtrToProtoTimestamp(fp.UpdatedAt),
		}
		protoParts = append(protoParts, protoPart)
	}
	response.GetPartsResponse = protoParts
	return response, nil
}

func (s *InventoryService) CreatePart(_ context.Context, r *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	if err := r.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	newUUID := uuid.New()
	part := &models.Part{
		ID:            newUUID,
		Name:          r.GetName(),
		Description:   r.GetDescription(),
		Price:         r.GetPrice(),
		StockQuantity: r.GetStockQuantity(),
		Category:      CategoryFromProtoToModel(r.GetCategory()),
		Dimensions:    DimensionsFromProtoToModel(r.GetDimensions()),
		Manufacturer:  ManufacturerFromProtoToModel(r.GetManufacturer()),
		Tags:          r.GetTags(),
		Metadata:      MetadataFromProtoToModel(r.GetMetadata()),
		CreatedAt:     timePtr(time.Now()),
		UpdatedAt:     nil,
	}
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.storage[newUUID] = part
	return &inventoryV1.CreatePartResponse{
		Uuid:          part.ID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryFromModelToProto(part.Category),
		Dimensions:    DimensionsFromModelToProto(part.Dimensions),
		Manufacturer:  ManufacturerFromModelToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataFromModelToProto(part.Metadata),
		CreatedAt:     TimePtrToProtoTimestamp(part.CreatedAt),
	}, nil
}

func MetadataFromProtoToModel(metadata map[string]*inventoryV1.Value) map[string]*models.Value {
	res := make(map[string]*models.Value, len(metadata))

	for k, v := range metadata {
		res[k] = &models.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}

	return res
}

func MetadataFromModelToProto(metadata map[string]*models.Value) map[string]*inventoryV1.Value {
	res := make(map[string]*inventoryV1.Value, len(metadata))

	for k, v := range metadata {
		res[k] = &inventoryV1.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}

	return res
}

func ManufacturerFromProtoToModel(manufacturer *inventoryV1.Manufacturer) *models.Manufacturer {
	return &models.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ManufacturerFromModelToProto(manufacturer *models.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func DimensionsFromProtoToModel(dimensions *inventoryV1.Dimensions) *models.Dimensions {
	return &models.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func DimensionsFromModelToProto(dimensions *models.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func CategoryFromProtoToModel(cat inventoryV1.Category) models.Category {
	switch cat {
	case inventoryV1.Category_CATEGORY_UNKNOWN:
		return models.CATEGORY_UNKNOWN
	case inventoryV1.Category_CATEGORY_ENGINE:
		return models.CATEGORY_ENGINE
	case inventoryV1.Category_CATEGORY_FUEL:
		return models.CATEGORY_FUEL
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return models.CATEGORY_PORTHOLE
	case inventoryV1.Category_CATEGORY_WING:
		return models.CATEGORY_WING
	default:
		return models.CATEGORY_UNKNOWN
	}
}

func CategoryFromModelToProto(cat models.Category) inventoryV1.Category {
	switch cat {
	case models.CATEGORY_UNKNOWN:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	case models.CATEGORY_ENGINE:
		return inventoryV1.Category_CATEGORY_ENGINE
	case models.CATEGORY_FUEL:
		return inventoryV1.Category_CATEGORY_FUEL
	case models.CATEGORY_PORTHOLE:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case models.CATEGORY_WING:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// TimePtrToProtoTimestamp конвертирует *time.Time в protobuf timestamp
func TimePtrToProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}
