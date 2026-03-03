package main

import (
	"context"
	"encoding/json"
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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Order struct {
	OrderID  uuid.UUID `json:"id"`
	Quantity int64     `json:"quantity"`
}

type Cache struct {
	storage map[uuid.UUID]*Order
	mu      *sync.RWMutex
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
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	storage := NewCache()
	r.Route("/api/v1/orders", func(r chi.Router) {
		r.Get("/{id}", GetOrder(storage))
		r.Post("/", CreateOrder(storage))
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", fmt.Sprintf("%d", port)),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Starting server oh port %d", port)
		err := server.ListenAndServe()
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
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server")
	}

	log.Printf("Server was stopped successfully")
}

func GetOrder(cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, urlParamOrderID)
		cache.mu.RLock()
		defer cache.mu.RUnlock()
		orderUUID, err := uuid.Parse(orderID)
		if err != nil {
			http.Error(w, "invalid order id", http.StatusBadRequest)
			return
		}
		order, ok := cache.storage[orderUUID]
		if !ok {
			http.Error(w, fmt.Sprintf("order with id %s not found", orderUUID), http.StatusNotFound)
			return
		}
		render.JSON(w, r, order)
		render.Status(r, http.StatusOK)
	}
}

func CreateOrder(cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order Order
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&order); err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}
		cache.mu.Lock()
		defer cache.mu.Unlock()
		newUUID := uuid.New()
		order.OrderID = newUUID
		cache.storage[newUUID] = &order

		w.Header().Set("Location", fmt.Sprintf(r.RequestURI+"/%s", newUUID))
		render.JSON(w, r, order)
		render.Status(r, http.StatusCreated)
	}
}
