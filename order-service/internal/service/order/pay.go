package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	"github.com/defan6/space-app/order-service/internal/service"
	"github.com/google/uuid"
)

func (dos *defaultOrderService) PayOrder(ctx context.Context, uuid uuid.UUID, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	op := "order-service#PayOrder"
	exists, err := dos.repo.Exists(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if !exists {
		return nil, fmt.Errorf("%s:%w", op, service.ErrNotFound)
	}

	order := converter.ConvertFromPayOrderRequestToRepoOrder(req)

	updatedOrder, err := dos.repo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.ConvertFromRepoOrderToPayOrderResponse(updatedOrder), nil
}
