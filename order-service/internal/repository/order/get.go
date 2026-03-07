package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	"github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

func (imr *inMemRepo) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.GetOrderResponse, error) {
	op := "order-repo#GetOrder"
	order, err := imr.getOrder(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return converter.ConvertFromRepoOrderToGetOrderResponse(order), nil
}

func (imr *inMemRepo) getOrder(_ context.Context, uuid uuid.UUID) (*repomodel.Order, error) {
	op := "storage#GetOrder"
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	res, ok := imr.cache[uuid]
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, repomodel.ErrNotFound)
	}

	return res, nil
}
