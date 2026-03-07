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

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Parts           []string
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
