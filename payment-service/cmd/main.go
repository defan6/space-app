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

	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	grpcPort        = 50051
	httpPort        = 8082
	shutdownTimeout = time.Second * 5
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) PayOrder(_ context.Context, r *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := r.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	trUUID := uuid.New().String()
	response := &paymentV1.PayOrderResponse{
		TransactionUuid: trUUID,
	}
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		panic("cannot start listening port")
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Cannot close listening on port %d\n", grpcPort)
		}
	}()

	s := grpc.NewServer()

	service := NewPaymentService()

	paymentV1.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf(" grpc server started on port: %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve grpc server on port: %d\n", grpcPort)
		}
	}()
	var gatewayServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		fileServer := http.FileServer(http.Dir("../shared/api/payment/v1/swagger"))

		httpMux := http.NewServeMux()

		httpMux.Handle("/api/v1/payment", mux)
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/payment.swagger.json", fileServer)

		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}

			fileServer.ServeHTTP(w, r)
		}))

		if err != nil {
			log.Printf("Failed to register grpc gateway: %v\n", err)
			return
		}

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
	log.Printf("Stopping servers...\n")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if gatewayServer != nil {
			err = gatewayServer.Shutdown(ctx)
			if err != nil {
				log.Printf("failed to shutdown http server on port %d, %v\n", httpPort, err)
			}
		}
		s.GracefulStop()
	}()

	select {
	case <-done:
		log.Printf("Planning stopped servers\n")
	case <-ctx.Done():
		log.Printf("Context deadline exceeded. Terminate stop\n")
	}
}
