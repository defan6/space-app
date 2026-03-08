package order

import (
	"context"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

// GetOrders implements GetOrders operation.
//
// Get Orders.
//
// GET /api/v1/orders
func (oa *orderAPI) GetOrders(ctx context.Context) (orderV1.GetOrdersRes, error) {

	orders, err := oa.service.GetOrders(ctx)
	if err != nil {
		return &orderV1.InternalServerError{
			Message:   "internal server error",
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}

	return converter.FromServiceGetOrdersResponse(orders), nil
}
