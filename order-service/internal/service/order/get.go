package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	"github.com/google/uuid"
)

func (dos *defaultOrderService) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.GetOrderResponse, error) {
	op := "order-service#GetOrder"
	order, err := dos.repo.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return converter.ConvertFromRepoOrderToGetOrderResponse(order), nil
}
