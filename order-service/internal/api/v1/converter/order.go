package converter

import (
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
)

func ConvertFromAPIStatusToStatus(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCancelled
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPaid
	case orderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPendingPayment
	default:
		return model.OrderStatusCancelled
	}
}

func ConvertFromStatusToAPIStatus(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	default:
		return orderV1.OrderStatusCANCELLED
	}
}

func ConvertFromAPIPaymentMethodToPaymentMethod(pm orderV1.PaymentMethod) (model.PaymentMethod, error) {
	switch pm {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCredit, nil
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnknown, nil
	}
}

func ConvertFromPaymentMethodToAPIPaymentMethod(pm model.PaymentMethod) orderV1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodCredit:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
func FromAPICreateOrderRequest(request *orderV1.CreateOrderRequest) (*model.CreateOrderRequest, error) {
	op := "order-api-converter#FromAPICreateOrderRequest"
	parts, err := FromAPIPartItems(request.PartItems)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &model.CreateOrderRequest{
		UserUUID:  request.UserUUID,
		PartItems: parts,
	}, nil
}

func FromAPIPartItems(items []orderV1.PartItem) ([]model.PartItem, error) {
	parts := make([]model.PartItem, 0, len(items))
	op := "order-api-converter#FromAPIPartItems"
	for _, p := range items {
		part, err := FromAPIPartItem(p)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func FromAPIPartItem(item orderV1.PartItem) (model.PartItem, error) {
	return model.PartItem{
		PartUUID: item.GetPartUUID(),
		Quantity: item.Quantity,
		Price:    item.Price,
	}, nil
}

func FromServicePartItems(items []model.PartItem) []orderV1.PartItem {
	parts := make([]orderV1.PartItem, 0, len(items))
	for _, p := range items {
		part := FromServicePartItem(p)
		parts = append(parts, part)
	}

	return parts
}

func FromServicePartItem(item model.PartItem) orderV1.PartItem {
	return orderV1.PartItem{
		PartUUID: item.PartUUID,
		Quantity: item.Quantity,
		Price:    item.Price,
	}
}

func FromServiceCreateOrderResponse(response *model.CreateOrderResponse) *orderV1.CreateOrderResponse {
	return &orderV1.CreateOrderResponse{
		OrderUUID:  response.OrderUUID,
		TotalPrice: response.TotalPrice,
	}
}

func FromServiceGetOrderResponse(o *model.GetOrderResponse) *orderV1.GetOrderResponse {
	return &orderV1.GetOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartItems:       FromServicePartItems(o.PartItems),
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   ConvertFromPaymentMethodToAPIPaymentMethod(o.PaymentMethod),
		TransactionUUID: o.TransactionUUID,
		Status:          ConvertFromStatusToAPIStatus(o.Status),
	}
}

func FromServiceGetOrdersResponse(response *model.GetOrdersResponse) []orderV1.GetOrderResponse {
	apiOrders := make([]orderV1.GetOrderResponse, 0, len(response.Orders))

	for _, ro := range response.Orders {
		o := FromServiceGetOrderResponse(ro)
		apiOrders = append(apiOrders, *o)
	}
	return apiOrders
}

func FromAPIPayOrderRequest(request *orderV1.PayOrderRequest) (*model.PayOrderRequest, error) {
	op := "order-api-converter#FromAPIPayOrderRequest"
	pm, err := ConvertFromAPIPaymentMethodToPaymentMethod(request.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &model.PayOrderRequest{
		UserUUID:      request.UserUUID,
		PaymentMethod: pm,
	}, nil
}

func FromServicePayOrderResponse(response *model.PayOrderResponse) *orderV1.PayOrderResponse {
	return &orderV1.PayOrderResponse{
		TransactionUUID: response.TransactionUUID,
	}
}

func FromServiceCancelOrderResponse(response *model.CancelOrderResponse) *orderV1.CancelOrderResponse {
	return &orderV1.CancelOrderResponse{
		OrderUUID: response.OrderUUID,
		Status:    ConvertFromStatusToAPIStatus(response.Status),
	}
}
