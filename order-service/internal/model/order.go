package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserUUID  uuid.UUID
	PartItems []PartItemRequest
}

type CreateOrderResponse struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	Parts           []PartItemResponse
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
	PartItems       []PartItemResponse
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUUID uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
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

type CheckInventoryResponse struct {
	PartItems  []PartInfo
	TotalPrice float64
}
