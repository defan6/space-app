package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

// PayOrder implements PayOrder operation.
//
// Pay Order.
//
// POST /api/v1/orders/{uuid}/pay
func (oa *orderAPI) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	op := "order-service#PayOrder"
	ouuid, err := uuid.Parse(params.UUID)
	if err != nil {
		return &orderV1.BadRequestError{
			Message:   fmt.Sprintf("failed to parse uuid: %s", params.UUID),
			ErrorCode: "BAD_REQUEST",
		}, nil
	}

	serviceReq, err := converter.FromAPIPayOrderRequest(req)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	payOrderRes, err := oa.service.PayOrder(ctx, ouuid, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.FromServicePayOrderResponse(payOrderRes), nil
}
