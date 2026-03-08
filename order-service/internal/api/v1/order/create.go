package apiorder

import (
	"context"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

// CreateOrder implements CreateOrder operation.
//
// Create Order.
//
// POST /api/v1/orders
func (oa *orderAPI) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	serviceReq, err := converter.FromAPICreateOrderRequest(req)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   "failed to create order",
			ErrorCode: "BAD_REQUEST",
		}, nil
	}

	createdOrder, err := oa.service.CreateOrder(ctx, serviceReq)
	if err != nil {
		return &orderV1.InternalServerError{
			Message:   "failed to create order",
			ErrorCode: "INTERNAL_SERVER_ERROR",
		}, nil
	}

	apiRes := converter.FromServiceCreateOrderResponse(createdOrder)
	return apiRes, nil
}
