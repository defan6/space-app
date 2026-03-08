package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
)

func (dos *defaultOrderService) GetOrders(ctx context.Context) (*model.GetOrdersResponse, error) {
	const op = "order-service#GetOrders"
	orders, err := dos.repo.GetOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return converter.ConvertFromRepoOrdersToGetOrdersResponse(orders), nil
}
