package payment

import (
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type paymentClient struct {
	paymentV1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *paymentClient {
	return &paymentClient{
		PaymentServiceClient: paymentV1.NewPaymentServiceClient(conn),
	}
}
