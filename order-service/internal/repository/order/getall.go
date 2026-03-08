package repoorder

import (
	"context"
	"log"

	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
)

func (imr *inMemRepo) GetOrders(ctx context.Context) ([]*repomodel.Order, error) {
	op := "order-repo#GetOrders"
	orders := make([]*repomodel.Order, 0, 10)

	imr.mu.RLock()
	for _, o := range imr.cache {
		orders = append(orders, o)
	}
	imr.mu.RUnlock()
	log.Printf("%s: %#v", op, orders)
	return orders, nil
}
