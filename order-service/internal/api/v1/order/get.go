package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	"github.com/defan6/space-app/order-service/internal/service"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

// GetOrder implements GetOrder operation.
//
// Get Order.
//
// GET /api/v1/orders/{uuid}
func (oa *orderAPI) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {

	order, err := oa.service.GetOrder(ctx, params.UUID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &orderV1.NotFoundError{
				Message:   fmt.Sprintf("order with id %s not found", params.UUID.String()),
				ErrorCode: "ORDER_NOT_FOUND",
			}, nil
		}
		return &orderV1.InternalServerError{
			Message:   "internal server error",
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}

	return converter.FromServiceGetOrderResponse(order), nil
}
