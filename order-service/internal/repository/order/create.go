package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

func (imr *inMemRepo) CreateOrder(ctx context.Context, order *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	op := "order-repo#CreateOrder"
	repoOrder := converter.ConvertFromCreateOrderRequestToRepoOrder(order)
	savedOrder, err := imr.createOrder(ctx, repoOrder)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.ConvertFromOrderToCreateOrderResponse(savedOrder), nil
}

func (imr *inMemRepo) createOrder(_ context.Context, order *repomodel.Order) (*repomodel.Order, error) {
	ouuid := uuid.New()
	order.OrderUUID = ouuid
	imr.mu.Lock()
	imr.cache[ouuid] = order
	imr.mu.Unlock()
	return order, nil
}
