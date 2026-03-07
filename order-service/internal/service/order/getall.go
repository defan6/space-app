package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository"
)

type defaultOrderService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *defaultOrderService {
	return &defaultOrderService{
		repo: repo,
	}
}

func (dos *defaultOrderService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	const op = "order-service#GetOrders"
	orders, err := dos.repo.GetOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return orders, nil
}
