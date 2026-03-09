package model

import (
	"time"

	"github.com/google/uuid"
)

// PartItemRequest — запрос на получение информации о детали.
// Используется в CreateOrderRequest для указания количества.
type PartItemRequest struct {
	PartUUID uuid.UUID
	Quantity int64
}

// PartInfo — минимальная информация о детали для внутренней логики.
// Используется при проверке inventory и расчёте стоимости.
type PartInfo struct {
	PartUUID uuid.UUID
	Quantity int64
	Price    float64
}

// PartItem — полная информация о детали из inventory service.
// Используется внутри сервиса для работы с данными о деталях.
type PartItem struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// PartItemResponse — деталь в ответе API.
// Используется в CreateOrderResponse и GetOrderResponse.
type PartItemResponse struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
