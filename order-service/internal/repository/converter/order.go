package converter

import (
	"github.com/defan6/space-app/order-service/internal/model"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
)

func ConvertFromRepoOrdersToGetOrdersResponse(repoOrders []*repomodel.Order) *model.GetOrdersResponse {
	orders := make([]*model.GetOrderResponse, 0, len(repoOrders))

	for _, ro := range repoOrders {
		o := ConvertFromRepoOrderToGetOrderResponse(ro)
		orders = append(orders, o)
	}
	return &model.GetOrdersResponse{
		Orders: orders,
	}
}

func ConvertFromRepoOrderToGetOrderResponse(ro *repomodel.Order) *model.GetOrderResponse {
	return &model.GetOrderResponse{
		OrderUUID:       ro.OrderUUID,
		UserUUID:        ro.UserUUID,
		Parts:           ro.Parts,
		TotalPrice:      ro.TotalPrice,
		PaymentMethod:   ConvertFromRepoPaymentMethodToPaymentMethod(ro.PaymentMethod),
		TransactionUUID: ro.TransactionUUID,
		Status:          ConvertFromRepoStatusToStatus(ro.Status),
	}
}

func ConvertFromRepoStatusToStatus(status repomodel.OrderStatus) model.OrderStatus {
	switch status {
	case repomodel.OrderStatusCancelled:
		return model.OrderStatusCancelled
	case repomodel.OrderStatusPaid:
		return model.OrderStatusPaid
	case repomodel.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment
	default:
		return model.OrderStatusCancelled
	}
}

func ConvertFromStatusToRepoStatus(status model.OrderStatus) repomodel.OrderStatus {
	switch status {
	case model.OrderStatusCancelled:
		return repomodel.OrderStatusCancelled
	case model.OrderStatusPaid:
		return repomodel.OrderStatusPaid
	case model.OrderStatusPendingPayment:
		return repomodel.OrderStatusPendingPayment
	default:
		return repomodel.OrderStatusCancelled
	}
}

func ConvertFromRepoPaymentMethodToPaymentMethod(pm repomodel.PaymentMethod) model.PaymentMethod {
	switch pm {
	case repomodel.PaymentMethodCard:
		return model.PaymentMethodCard
	case repomodel.PaymentMethodCredit:
		return model.PaymentMethodCredit
	case repomodel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repomodel.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func ConvertFromPaymentMethodToRepoPaymentMethod(pm model.PaymentMethod) repomodel.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return repomodel.PaymentMethodCard
	case model.PaymentMethodCredit:
		return repomodel.PaymentMethodCredit
	case model.PaymentMethodSBP:
		return repomodel.PaymentMethodSBP
	case model.PaymentMethodInvestorMoney:
		return repomodel.PaymentMethodInvestorMoney
	default:
		return repomodel.PaymentMethodUnknown
	}
}

func ConvertFromCreateOrderRequestToRepoOrder(o *model.CreateOrderRequest) *repomodel.Order {
	return &repomodel.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		Parts:           o.Parts,
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   ConvertFromPaymentMethodToRepoPaymentMethod(o.PaymentMethod),
		TransactionUUID: o.TransactionUUID,
		Status:          ConvertFromStatusToRepoStatus(o.Status),
	}
}

func ConvertFromPayOrderRequestToRepoOrder(o *model.PayOrderRequest) *repomodel.Order {
	return &repomodel.Order{
		OrderUUID:     o.OrderUUID,
		UserUUID:      o.UserUUID,
		PaymentMethod: ConvertFromPaymentMethodToRepoPaymentMethod(o.PaymentMethod),
	}
}

func ConvertFromRepoOrderToPayOrderResponse(o *repomodel.Order) *model.PayOrderResponse {
	return &model.PayOrderResponse{
		OrderUUID:     o.OrderUUID,
		UserUUID:      o.UserUUID,
		PaymentMethod: ConvertFromRepoPaymentMethodToPaymentMethod(o.PaymentMethod),
	}
}

func ConvertFromOrderToCreateOrderResponse(o *repomodel.Order) *model.CreateOrderResponse {
	return &model.CreateOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		Parts:           o.Parts,
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   ConvertFromRepoPaymentMethodToPaymentMethod(o.PaymentMethod),
		TransactionUUID: o.TransactionUUID,
		Status:          ConvertFromRepoStatusToStatus(o.Status),
	}
}
