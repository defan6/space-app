package model

// OrderStatus — статус заказа.
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты.
type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCredit        PaymentMethod = "CREDIT"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Category — категория детали.
type Category string

const (
	CategoryUnknown  Category = "CATEGORY_UNKNOWN"
	CategoryEngine   Category = "CATEGORY_ENGINE"
	CategoryFuel     Category = "CATEGORY_FUEL"
	CategoryPorthole Category = "CATEGORY_PORTHOLE"
	CategoryWing     Category = "CATEGORY_WING"
)

// Dimensions — габариты детали.
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer — производитель детали.
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Value — мета-данные (значение в формате key-value).
// Хранит значение в типизированном виде.
type Value struct {
	StringValue string
	Int64Value  int64
	DoubleValue float64
	BoolValue   bool
}
