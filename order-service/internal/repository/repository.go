package repository

import (
	"context"

	"github.com/defan6/space-app/order-service/internal/model"
)

type Repository interface {
	CancelOrder(ctx context.Context, uuid string) (*model.CancelOrderResponse, error)
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	GetOrder(ctx context.Context, uuid string) (*model.GetOrderResponse, error)
	GetOrders(ctx context.Context) ([]*model.GetOrdersResponse, error)
	PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error)
}
