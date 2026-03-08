package repomodel

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

type PartItem struct {
	PartUUID uuid.UUID
	Quantity int64
	Price    float64
}
type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartsItems      []PartItem
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
