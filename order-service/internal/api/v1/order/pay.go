package apiorder

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/api/v1/converter"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

// PayOrder implements PayOrder operation.
//
// Pay Order.
//
// POST /api/v1/orders/{uuid}/pay
func (oa *orderAPI) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	op := "order-service#PayOrder"

	serviceReq, err := converter.FromAPIPayOrderRequest(req)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	payOrderRes, err := oa.service.PayOrder(ctx, req.UserUUID, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.FromServicePayOrderResponse(payOrderRes), nil
}
