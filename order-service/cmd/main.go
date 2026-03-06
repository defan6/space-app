package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Parts           []string
	TotalPrice      float64
	PaymentMethod   orderV1.PaymentMethod
	TransactionUUID uuid.UUID
	Status          orderV1.OrderStatus
}

type Cache struct {
	storage map[uuid.UUID]*Order
	mu      *sync.RWMutex
}

type PaymentClient struct {
	client paymentV1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	client := paymentV1.NewPaymentServiceClient(conn)
	return &PaymentClient{
		client: client,
	}
}

type Handler struct {
	paymentClient *PaymentClient
	storage       *Cache
}

func NewHandler(storage *Cache, paymentClient *PaymentClient) *Handler {
	return &Handler{
		storage:       storage,
		paymentClient: paymentClient,
	}
}

func NewCache() *Cache {
	return &Cache{
		storage: make(map[uuid.UUID]*Order),
		mu:      &sync.RWMutex{},
	}
}

const (
	port                  = 8080
	readHeaderTimeout     = 10 * time.Second
	shutdownTimeout       = 5 * time.Second
	paymentServiceAddress = 50051
)

func main() {
	r := chi.NewRouter()
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", paymentServiceAddress),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create new grpc client on port: %d", paymentServiceAddress)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connection on port: %d", paymentServiceAddress)
		}
	}()
	cache := NewCache()
	paymentClient := NewPaymentClient(conn)

	h := NewHandler(cache, paymentClient)
	orderServer, err := orderV1.NewServer(h)
	if err != nil {
		panic("Failed to start order server")
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", orderServer)

	fileServer := http.FileServer(http.Dir("../shared/api/order/v1/swagger"))
	r.Handle("/swagger-ui.html", fileServer)
	r.Handle("/order.swagger.json", fileServer)

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		fileServer.ServeHTTP(w, r)
	}))
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", fmt.Sprintf("%d", port)),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Starting server oh port %d", port)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to start server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server")
	}

	log.Printf("Server was stopped successfully")
}

// CancelOrder implements CancelOrder operation.
//
// Cancel Order.
//
// DELETE /api/v1/orders/{uuid}/cancel
func (h *Handler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	id := params.UUID
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("invalid uuid: %s", id),
			ErrorCode: "INVALID_ORDER_UUID",
		}, nil
	}

	h.storage.mu.RLock()
	order, ok := h.storage.storage[uuidFromId]
	h.storage.mu.RUnlock()
	if !ok {
		return &orderV1.NotFoundError{
			Message:   fmt.Sprintf("order with uuid %s not found", uuidFromId),
			ErrorCode: "ORDER_NOT_FOUND",
		}, nil
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Message:   fmt.Sprintf("order with uuid %s already cancelled", uuidFromId),
			ErrorCode: "ORDER_ALREADY_CANCELLED",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Message:   fmt.Sprintf("order with uuid %s already payed", uuidFromId),
			ErrorCode: "ORDER_ALREADY_PAYED",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED
	h.storage.mu.Lock()
	h.storage.storage[uuidFromId] = order
	defer h.storage.mu.Unlock()
	return &orderV1.CancelOrderNoContent{}, nil
}

// CreateOrder implements CreateOrder operation.
//
// Create Order.
//
// POST /api/v1/orders
func (h *Handler) CreateOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	newUUID := uuid.New()
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()
	maxPrice := big.NewInt(9000)
	randomInt, err := rand.Int(rand.Reader, maxPrice)
	if err != nil {
		log.Printf("error generating int")
		return &orderV1.InternalServerError{
			Message:   "Internal Server Error",
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}
	totalPrice := float64(randomInt.Int64())
	userUUID, err := uuid.Parse(req.UserUUID.Value)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("Invalid user uuid: %s", req.UserUUID.Value),
			ErrorCode: "INVALID_USER_UUID",
		}, nil
	}
	order := &Order{
		UserUUID:      userUUID,
		OrderUUID:     newUUID,
		Parts:         req.GetPartUuids(),
		TotalPrice:    totalPrice,
		PaymentMethod: orderV1.PaymentMethodUNKNOWN,
		Status:        orderV1.OrderStatusPENDINGPAYMENT,
	}
	h.storage.storage[newUUID] = order
	response := &orderV1.CreateOrderResponse{
		OrderUUID:  orderV1.NewOptString(newUUID.String()),
		TotalPrice: orderV1.NewOptFloat64(totalPrice),
	}
	return response, nil
}

// GetOrder implements GetOrder operation.
//
// Get Order.
//
// GET /api/v1/orders/{uuid}
func (h *Handler) GetOrder(_ context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	h.storage.mu.RLock()
	defer h.storage.mu.RUnlock()
	id := params.UUID
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("invalid uuid: %s", id),
			ErrorCode: "INVALID_UUID",
		}, nil
	}
	order, ok := h.storage.storage[uuidFromId]
	if !ok {
		return &orderV1.NotFoundError{
			Message:   fmt.Sprintf("order with uuid %s not found", uuidFromId),
			ErrorCode: "ORDER_NOT_FOUND",
		}, nil
	}
	response := &orderV1.GetOrderResponse{
		OrderUUID:       orderV1.NewOptString(uuidFromId.String()),
		UserUUID:        orderV1.NewOptString(order.UserUUID.String()),
		PartUuids:       order.Parts,
		TotalPrice:      orderV1.NewOptFloat64(order.TotalPrice),
		TransactionUUID: orderV1.NewOptString(order.TransactionUUID.String()),
		PaymentMethod:   orderV1.NewOptPaymentMethod(order.PaymentMethod),
		Status:          orderV1.NewOptOrderStatus(order.Status),
	}
	return response, nil
}

// GetOrders implements GetOrders operation.
//
// Get Orders.
//
// GET /api/v1/orders
func (h *Handler) GetOrders(_ context.Context) (orderV1.GetOrdersRes, error) {
	response := &orderV1.GetOrdersResponse{}
	for _, o := range h.storage.storage {
		ores := orderV1.GetOrderResponse{
			OrderUUID:       orderV1.NewOptString(o.OrderUUID.String()),
			UserUUID:        orderV1.NewOptString(o.UserUUID.String()),
			PartUuids:       o.Parts,
			TotalPrice:      orderV1.NewOptFloat64(o.TotalPrice),
			TransactionUUID: orderV1.NewOptString(o.TransactionUUID.String()),
			PaymentMethod:   orderV1.NewOptPaymentMethod(o.PaymentMethod),
			Status:          orderV1.NewOptOrderStatus(o.Status),
		}
		*response = append(*response, ores)
	}
	return response, nil
}

// PayOrder implements PayOrder operation.
//
// Pay Order.
//
// POST /api/v1/orders/{uuid}/pay
func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderId := params.UUID
	orderUUIDFromID, err := uuid.Parse(orderId)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("invalid uuid: %s", orderId),
			ErrorCode: "INVALID_ORDER_UUID",
		}, nil
	}

	h.storage.mu.RLock()
	order, ok := h.storage.storage[orderUUIDFromID]
	h.storage.mu.RUnlock()
	if !ok {
		return &orderV1.NotFoundError{
			Message:   fmt.Sprintf("order with uuid %s not found", order.OrderUUID),
			ErrorCode: "ORDER_NOT_FOUND",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Message:   fmt.Sprintf("order with uuid %s already paid", orderUUIDFromID),
			ErrorCode: "ORDER_ALREADY_PAID",
		}, nil
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Message:   fmt.Sprintf("order with uuid %s already cancelled", orderUUIDFromID),
			ErrorCode: "ORDER_ALREADY_CANCELLED",
		}, nil
	}
	payReq := &paymentV1.PayOrderRequest{
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		OrderUuid:     orderUUIDFromID.String(),
		UserUuid:      req.UserUUID.Value,
	}
	payRes, err := h.paymentClient.client.PayOrder(ctx, payReq)
	if err != nil {
		log.Printf("failed to pay order: %s, %v", orderUUIDFromID, err)
		return &orderV1.InternalServerError{
			Message:   fmt.Sprintf("internal server error"),
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}
	trID := payRes.TransactionUuid
	trUUID, err := uuid.Parse(trID)
	if err != nil {
		log.Printf("failed to pay order: %s, %v", orderUUIDFromID, err)
		return &orderV1.InternalServerError{
			Message:   fmt.Sprintf("internal server error"),
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = trUUID
	order.PaymentMethod = req.PaymentMethod.Value
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()
	h.storage.storage[order.OrderUUID] = order
	response := &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.NewOptString(trUUID.String()),
	}
	return response, nil
}
