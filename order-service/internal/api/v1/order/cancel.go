package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	"github.com/defan6/space-app/order-service/internal/service"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

// CancelOrder implements CancelOrder operation.
//
// Cancel Order.
//
// DELETE /api/v1/orders/{uuid}/cancel
func (oa *orderAPI) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {

	ouuid := params.UUID
	cancelRes, err := oa.service.CancelOrder(ctx, ouuid)

	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &orderV1.NotFoundError{
				Message:   fmt.Sprintf("order with id %s not found", ouuid.String()),
				ErrorCode: "ORDER_NOT_FOUND",
			}, nil
		}

		return &orderV1.InternalServerError{
			Message:   "internal server error",
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}

	return converter.FromServiceCancelOrderResponse(cancelRes), nil
}
