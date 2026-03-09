package converter

import (
	"fmt"

	"github.com/defan6/space-app/order-service/internal/model"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	"github.com/samber/lo"
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
		return model.OrderStatusPendingPayment
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
		return orderV1.OrderStatusPENDINGPAYMENT
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

func FromAPIPartItems(items []orderV1.PartItemRequest) ([]model.PartItemRequest, error) {
	parts := make([]model.PartItemRequest, 0, len(items))
	op := "order-api-converter#FromAPIPartItems"
	for _, p := range items {
		part, err := FromAPIPartItemRequest(p)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func FromAPIPartItemRequest(item orderV1.PartItemRequest) (model.PartItemRequest, error) {
	return model.PartItemRequest{
		PartUUID: item.GetPartUUID(),
		Quantity: item.Quantity,
	}, nil
}

func FromServicePartItemsResponse(items []model.PartItemResponse) []orderV1.PartItemResponse {
	parts := make([]orderV1.PartItemResponse, 0, len(items))
	for _, p := range items {
		part := FromServicePartItemResponse(p)
		parts = append(parts, part)
	}

	return parts
}

func FromServiceGetOrderResponse(o *model.GetOrderResponse) *orderV1.GetOrderResponse {
	return &orderV1.GetOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartItems:       FromServicePartItemsResponse(o.PartItems),
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   ConvertFromPaymentMethodToAPIPaymentMethod(o.PaymentMethod),
		TransactionUUID: o.TransactionUUID,
		Status:          ConvertFromStatusToAPIStatus(o.Status),
	}
}

func FromServicePartItemResponse(item model.PartItemResponse) orderV1.PartItemResponse {
	return orderV1.PartItemResponse{
		PartUUID: item.ID,
		Quantity: item.StockQuantity,
		Price:    item.Price,
	}
}

func FromServiceCreateOrderResponse(response *model.CreateOrderResponse) *orderV1.CreateOrderResponse {
	return &orderV1.CreateOrderResponse{
		OrderUUID:  response.OrderUUID,
		TotalPrice: response.TotalPrice,
	}
}

func FromServiceGetOrdersResponse(response *model.GetOrdersResponse) *orderV1.GetOrdersResponse {
	apiOrders := make([]orderV1.GetOrderResponse, 0, len(response.Orders))

	for _, ro := range response.Orders {
		o := FromServiceGetOrderResponse(ro)
		apiOrders = append(apiOrders, *o)
	}
	return lo.ToPtr(orderV1.GetOrdersResponse(apiOrders))
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
