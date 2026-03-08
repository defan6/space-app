package model

import (
	"time"

	"github.com/google/uuid"
)

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

type CreateOrderRequest struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartItems       []PartItem
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
}

type PartItem struct {
	PartUUID uuid.UUID
	Quantity int64
	Price    float64
}

type CreateOrderResponse struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Parts           []PartItem
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type GetOrderResponse struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartItems       []PartItem
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type CheckInventoryResponse struct {
	PartItems  []PartItem
	TotalPrice float64
}

type GetOrdersResponse struct {
	Orders []*GetOrderResponse
}

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
	Status          OrderStatus
}

type CancelOrderResponse struct {
	OrderUUID uuid.UUID
	Status    OrderStatus
}
