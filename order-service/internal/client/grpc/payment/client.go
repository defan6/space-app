package payment

import (
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type grpcPaymentClient struct {
	client paymentV1.PaymentServiceClient
}

func NewGrpcPaymentClient(conn *grpc.ClientConn) *grpcPaymentClient {
	client := paymentV1.NewPaymentServiceClient(conn)
	return &grpcPaymentClient{
		client: client,
	}
}
