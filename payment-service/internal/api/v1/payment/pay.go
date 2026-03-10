package payment

import (
	"context"

	"github.com/defan6/space-app/payment-service/internal/api/v1/converter"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pa *paymentAPI) PayOrder(ctx context.Context, r *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := r.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	req := converter.PayOrderRequestFromProto(r)
	resp, err := pa.paymentService.PayOrder(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.PayOrderResponseToProto(resp), nil
}
