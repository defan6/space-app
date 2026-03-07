package order

import (
	"context"

	"github.com/defan6/space-app/order-service/internal/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, uuid string, req *model.Order) (*model.Order, error)
	CreateOrder(ctx context.Context, uuid string, req *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	GetOrders(ctx context.Context) ([]*model.Order, error)
	PayOrder(ctx context.Context, uuid string, req *model.Order) (*model.Order, error)
}
