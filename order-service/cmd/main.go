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
	"syscall"
	"time"

	apiorder "github.com/defan6/space-app/order-service/internal/api/v1/order"
	"github.com/defan6/space-app/order-service/internal/client/grpc/inventory"
	"github.com/defan6/space-app/order-service/internal/client/grpc/payment"
	"github.com/defan6/space-app/order-service/internal/repository/order"
	serviceorder "github.com/defan6/space-app/order-service/internal/service/order"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port                    = 8080
	readHeaderTimeout       = 10 * time.Second
	shutdownTimeout         = 5 * time.Second
	paymentServiceAddress   = 50051
	inventoryServiceAddress = 50052
)

func main() {
	r := chi.NewRouter()
	pconn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", paymentServiceAddress),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create new grpc onn on port: %d", paymentServiceAddress)
	}
	defer func() {
		if cerr := pconn.Close(); cerr != nil {
			log.Printf("failed to close connection on port: %d", paymentServiceAddress)
		}
	}()
	iconn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", inventoryServiceAddress),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create new grpc conn on port: %d", inventoryServiceAddress)
	}
	defer func() {
		if cerr := iconn.Close(); cerr != nil {
			log.Printf("failed to close connection on port: %d", inventoryServiceAddress)
		}
	}()
	repo := repoorder.NewInMemRepo()
	paymentClient := payment.NewPaymentClient(pconn)
	inventoryClient := inventory.NewInventoryClient(iconn)
	service := serviceorder.NewService(repo, paymentClient, inventoryClient)
	api := apiorder.NewOrderHandler(service)
	orderServer, err := orderV1.NewServer(api)
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
