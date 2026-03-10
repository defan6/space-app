package payment

import (
	"context"
	"time"

	"github.com/defan6/space-app/payment-service/internal/model"
	"github.com/google/uuid"
)

func (dps *defaultPaymentService) PayOrder(ctx context.Context, r *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	trUUID := uuid.New()
	response := &model.PayOrderResponse{
		OrderUUID:       r.OrderUUID,
		UserUUID:        r.UserUUID,
		TransactionUUID: trUUID,
		PaymentMethod:   r.PaymentMethod,
	}
	return response, nil
}
