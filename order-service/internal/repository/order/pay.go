package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
)

func (imr *inMemRepo) PayOrder(ctx context.Context, order *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	op := "order-repo#PayOrder"
	repoOrder := converter.ConvertFromPayOrderRequestToRepoOrder(order)

	existingOrder, ok := imr.cache[order.OrderUUID]
	if !ok {
		imr.mu.RUnlock()
		return nil, repomodel.ErrNotFound
	}
	payedOrder, err := imr.payOrder(ctx, repoOrder)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.ConvertFromRepoOrderToPayOrderResponse(payedOrder), nil
}

func (imr *inMemRepo) payOrder(_ context.Context, order *repomodel.Order) (*repomodel.Order, error) {
	imr.mu.RLock()
	existingOrder, ok := imr.cache[order.OrderUUID]
	if !ok {
		imr.mu.RUnlock()
		return nil, repomodel.ErrNotFound
	}
	existingOrder.
		imr.mu.Lock()
	defer imr.mu.Unlock()
	imr.cache[]
}
