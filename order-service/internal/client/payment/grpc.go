package payment

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/client/converter"
	paymentmodel "github.com/defan6/space-app/order-service/internal/client/payment/model"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func (pc *paymentClient) ExternalPayOrder(ctx context.Context, orderUUID uuid.UUID, req *paymentV1.PayOrderRequest) (*paymentmodel.PayOrderInfo, error) {
	op := "payment-client#ExternalPayOrder"

	payReq := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      req.UserUuid,
		PaymentMethod: req.PaymentMethod,
	}

	payRes, err := pc.PayOrder(ctx, payReq)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	payInfo, err := converter.FromPaymentExternalPayOrderResponse(payRes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return payInfo, nil
}
