package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort        = 50051
	shutdownTimeout = time.Second * 5
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) PayOrder(context.Context, *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
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
		log.Printf("Server started on port: %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("Failed to serve on port: %d\n", grpcPort)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Stopping grpc server on port %d\n", grpcPort)
	timer := time.NewTimer(3 * time.Second)
	<-timer.C
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("Planning stopped server\n")
	case <-ctx.Done():
		log.Printf("Context deadline exceeded. Terminate stop\n")
	}
}
