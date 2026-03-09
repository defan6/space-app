package converter

import (
	"github.com/defan6/space-app/order-service/internal/model"
	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
)

// ConvertFromRepoOrdersToGetOrdersResponse конвертирует []repo.Order → model.GetOrdersResponse.
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

// ConvertFromRepoOrderToGetOrderResponse конвертирует repo.Order → model.GetOrderResponse.
func ConvertFromRepoOrderToGetOrderResponse(ro *repomodel.Order) *model.GetOrderResponse {
	return &model.GetOrderResponse{
		OrderUUID:       ro.OrderUUID,
		UserUUID:        ro.UserUUID,
		PartItems:       ConvertFromRepoPartsToPartResponses(ro.PartsItems),
		TotalPrice:      ro.TotalPrice,
		PaymentMethod:   ConvertFromRepoPaymentMethodToPaymentMethod(ro.PaymentMethod),
		TransactionUUID: ro.TransactionUUID,
		Status:          ConvertFromRepoStatusToStatus(ro.Status),
		CreatedAt:       ro.CreatedAt,
		UpdatedAt:       ro.UpdatedAt,
	}
}

// ConvertFromCreateOrderRequestToRepoOrder конвертирует model.CreateOrderRequest → repo.Order.
func ConvertFromCreateOrderRequestToRepoOrder(o *model.CreateOrderRequest) *repomodel.Order {
	return &repomodel.Order{
		UserUUID:   o.UserUUID,
		PartsItems: ConvertFromPartRequestsToRepoParts(o.PartItems),
	}
}

// ConvertFromPartsToRepoParts конвертирует []model.PartItem → []repo.PartItem.
func ConvertFromPartsToRepoParts(parts []model.PartItem) []repomodel.PartItem {
	res := make([]repomodel.PartItem, 0, len(parts))

	for _, part := range parts {
		repoPart := repomodel.PartItem{
			PartUUID: part.ID,
			Quantity: part.StockQuantity,
			Price:    part.Price,
		}
		res = append(res, repoPart)
	}

	return res
}

// ConvertFromPartRequestsToRepoParts конвертирует []model.PartItemRequest → []repo.PartItem.
func ConvertFromPartRequestsToRepoParts(parts []model.PartItemRequest) []repomodel.PartItem {
	res := make([]repomodel.PartItem, 0, len(parts))

	for _, part := range parts {
		repoPart := repomodel.PartItem{
			PartUUID: part.PartUUID,
			Quantity: part.Quantity,
			Price:    0, // Цена будет заполнена при проверке inventory
		}
		res = append(res, repoPart)
	}

	return res
}

// ConvertFromPartInfosToRepoParts конвертирует []model.PartInfo → []repo.PartItem.
func ConvertFromPartInfosToRepoParts(parts []model.PartInfo) []repomodel.PartItem {
	res := make([]repomodel.PartItem, 0, len(parts))

	for _, part := range parts {
		repoPart := repomodel.PartItem{
			PartUUID: part.PartUUID,
			Quantity: part.Quantity,
			Price:    part.Price,
		}
		res = append(res, repoPart)
	}

	return res
}

// ConvertFromRepoPartsToPartResponses конвертирует []repo.PartItem → []model.PartItemResponse.
func ConvertFromRepoPartsToPartResponses(parts []repomodel.PartItem) []model.PartItemResponse {
	res := make([]model.PartItemResponse, 0, len(parts))

	for _, part := range parts {
		repoPart := model.PartItemResponse{
			ID:            part.PartUUID,
			StockQuantity: part.Quantity,
			Price:         part.Price,
		}
		res = append(res, repoPart)
	}

	return res
}

// ConvertFromRepoStatusToStatus конвертирует repo.OrderStatus → model.OrderStatus.
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

// ConvertFromStatusToRepoStatus конвертирует model.OrderStatus → repo.OrderStatus.
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

// ConvertFromRepoPaymentMethodToPaymentMethod конвертирует repo.PaymentMethod → model.PaymentMethod.
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

// ConvertFromPaymentMethodToRepoPaymentMethod конвертирует model.PaymentMethod → repo.PaymentMethod.
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

// ConvertFromPayOrderRequestToRepoOrder конвертирует model.PayOrderRequest → repo.Order.
func ConvertFromPayOrderRequestToRepoOrder(o *model.PayOrderRequest) *repomodel.Order {
	return &repomodel.Order{
		OrderUUID:     o.OrderUUID,
		UserUUID:      o.UserUUID,
		PaymentMethod: ConvertFromPaymentMethodToRepoPaymentMethod(o.PaymentMethod),
	}
}

// ConvertFromRepoOrderToPayOrderResponse конвертирует repo.Order → model.PayOrderResponse.
func ConvertFromRepoOrderToPayOrderResponse(o *repomodel.Order) *model.PayOrderResponse {
	return &model.PayOrderResponse{
		OrderUUID:     o.OrderUUID,
		UserUUID:      o.UserUUID,
		PaymentMethod: ConvertFromRepoPaymentMethodToPaymentMethod(o.PaymentMethod),
	}
}

// ConvertFromOrderToCreateOrderResponse конвертирует repo.Order → model.CreateOrderResponse.
func ConvertFromOrderToCreateOrderResponse(o *repomodel.Order) *model.CreateOrderResponse {
	return &model.CreateOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		Parts:           ConvertFromRepoPartsToPartResponses(o.PartsItems),
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   ConvertFromRepoPaymentMethodToPaymentMethod(o.PaymentMethod),
		TransactionUUID: o.TransactionUUID,
		Status:          ConvertFromRepoStatusToStatus(o.Status),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

// ConvertFromOrderToCancelOrderResponse конвертирует repo.Order → model.CancelOrderResponse.
func ConvertFromOrderToCancelOrderResponse(o *repomodel.Order) *model.CancelOrderResponse {
	return &model.CancelOrderResponse{
		OrderUUID: o.OrderUUID,
		Status:    ConvertFromRepoStatusToStatus(o.Status),
	}
}
