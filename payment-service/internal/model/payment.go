package model

import "github.com/google/uuid"

type OrderStatus string
type PaymentMethod string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCredit        PaymentMethod = "CREDIT"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PayOrderRequest struct {
	OrderUUID     uuid.UUID
	UserUUID      uuid.UUID
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
}
