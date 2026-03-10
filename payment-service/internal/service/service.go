package service

import (
	"context"

	"github.com/defan6/space-app/payment-service/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, r *model.PayOrderRequest) (*model.PayOrderResponse, error)
}
