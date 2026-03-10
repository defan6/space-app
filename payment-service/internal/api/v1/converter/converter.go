package converter

import (
	"github.com/defan6/space-app/payment-service/internal/model"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

// PaymentMethodToProto конвертирует модель PaymentMethod в proto PaymentMethod
func PaymentMethodToProto(m model.PaymentMethod) paymentV1.PaymentMethod {
	switch m {
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCredit:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

// PaymentMethodFromProto конвертирует proto PaymentMethod в модель PaymentMethod
func PaymentMethodFromProto(m paymentV1.PaymentMethod) model.PaymentMethod {
	switch m {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCredit
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

// PayOrderRequestFromProto конвертирует proto PayOrderRequest в модель PayOrderRequest
func PayOrderRequestFromProto(r *paymentV1.PayOrderRequest) *model.PayOrderRequest {
	return &model.PayOrderRequest{
		OrderUUID:     uuid.MustParse(r.OrderUuid),
		UserUUID:      uuid.MustParse(r.UserUuid),
		PaymentMethod: PaymentMethodFromProto(r.PaymentMethod),
	}
}

// PayOrderResponseToProto конвертирует модель PayOrderResponse в proto PayOrderResponse
func PayOrderResponseToProto(r *model.PayOrderResponse) *paymentV1.PayOrderResponse {
	return &paymentV1.PayOrderResponse{
		TransactionUuid: r.TransactionUUID.String(),
	}
}
