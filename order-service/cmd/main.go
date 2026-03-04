package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
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

type Handler struct {
	storage *Cache
}

func NewHandler(storage *Cache) *Handler {
	return &Handler{
		storage: storage,
	}
}

func NewCache() *Cache {
	return &Cache{
		storage: make(map[uuid.UUID]*Order),
		mu:      &sync.RWMutex{},
	}
}

const (
	urlParamOrderID   = "id"
	port              = 8080
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

func main() {
	r := chi.NewRouter()
	cache := NewCache()
	h := NewHandler(cache)
	orderOpenAPIServer, err := orderV1.NewServer(h)
	if err != nil {
		panic("Failed to start order OPEN API server")
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", orderOpenAPIServer)

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
func (h *Handler) CancelOrder(ctx context.Context, req *orderV1.CancelOrderRequest, params orderV1.CancelOrderParams) (r orderV1.CancelOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateOrder implements CreateOrder operation.
//
// Create Order.
//
// POST /api/v1/orders
func (h *Handler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	newUUID := uuid.New()
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()
	totalPrice := rand.Float64() * float64(rand.Intn(9000))
	userUUID, err := uuid.Parse(req.UserUUID.Value)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("Invalid user uuid: %s", userUUID),
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
func (h *Handler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	h.storage.mu.RLock()
	defer h.storage.mu.RUnlock()
	id := params.UUID
	uuidFromId, err := uuid.Parse(id)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("invalid UUID: %s", uuidFromId),
			ErrorCode: "INVALID_UUID",
		}, nil
	}
	order, ok := h.storage.storage[uuidFromId]
	if !ok {
		return &orderV1.NotFoundError{
			Message:   fmt.Sprintf("order with id %s not found", uuidFromId),
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
func (h *Handler) GetOrders(ctx context.Context) (r orderV1.GetOrdersRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PayOrder implements PayOrder operation.
//
// Pay Order.
//
// POST /api/v1/orders/{uuid}/pay
func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (r orderV1.PayOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}
