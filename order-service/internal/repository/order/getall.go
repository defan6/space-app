package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
)

func (imr *inMemRepo) GetOrders(ctx context.Context) (*model.GetOrdersResponse, error) {
	op := "order-repo#GetOrders"
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	orders, err := imr.getOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.ConvertFromRepoOrdersToGetOrdersResponse(orders), nil
}

func (imr *inMemRepo) getOrders(_ context.Context) ([]*repomodel.Order, error) {
	orders := make([]*repomodel.Order, 0, 10)

	imr.mu.RLock()
	for _, o := range imr.cache {
		orders = append(orders, o)
	}
	imr.mu.RUnlock()

	return orders, nil
}
