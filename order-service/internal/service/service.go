package service

import (
	"context"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/google/uuid"
)

type OrderService interface {
	CancelOrder(ctx context.Context, uuid uuid.UUID) (*model.CancelOrderResponse, error)
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.GetOrderResponse, error)
	GetOrders(ctx context.Context) (*model.GetOrdersResponse, error)
	PayOrder(ctx context.Context, uuid uuid.UUID, req *model.PayOrderRequest) (*model.PayOrderResponse, error)
}
