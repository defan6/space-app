package serviceorder

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/defan6/space-app/order-service/internal/service"
	"github.com/google/uuid"
)

func (dos *defaultOrderService) CancelOrder(ctx context.Context, uuid uuid.UUID) (*model.CancelOrderResponse, error) {
	op := "order-service#CancelOrder"
	existingOrder, err := dos.repo.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, service.ErrNotFound)
	}

	if existingOrder.Status == repomodel.OrderStatusCancelled {
		return nil, fmt.Errorf("%s cannot cancel already cancelled order", op)
	}
	if existingOrder.Status == repomodel.OrderStatusPaid {
		return nil, fmt.Errorf("%s cannot cancel already paid order", op)
	}

	existingOrder.Status = repomodel.OrderStatusCancelled
	updated, err := dos.repo.Update(ctx, existingOrder)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return converter.ConvertFromOrderToCancelOrderResponse(updated), nil
}
