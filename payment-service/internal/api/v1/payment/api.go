package payment

import (
	"github.com/defan6/space-app/payment-service/internal/service"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
)

type paymentAPI struct {
	paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentAPI(service service.PaymentService) *paymentAPI {
	return &paymentAPI{
		paymentService: service,
	}
}
