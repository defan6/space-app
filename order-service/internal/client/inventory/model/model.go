package model

import (
	"time"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/google/uuid"
)

// PartItem — полная информация о детали из inventory service.
// Алиас на model.PartItem для удобства импорта в client слое.
//
// Deprecated: Используйте model.PartItem напрямую.
type PartItem = model.PartItem

// Category — категория детали.
//
// Deprecated: Используйте model.Category напрямую.
type Category = model.Category

// Dimensions — габариты детали.
//
// Deprecated: Используйте model.Dimensions напрямую.
type Dimensions = model.Dimensions

// Manufacturer — производитель детали.
//
// Deprecated: Используйте model.Manufacturer напрямую.
type Manufacturer = model.Manufacturer

// Value — мета-данные.
//
// Deprecated: Используйте model.Value напрямую.
type Value = model.Value

// Для совместимости со старым кодом
var (
	CategoryUnknown  = model.CategoryUnknown
	CategoryEngine   = model.CategoryEngine
	CategoryFuel     = model.CategoryFuel
	CategoryPorthole = model.CategoryPorthole
	CategoryWing     = model.CategoryWing
)

// Helper для создания PartItem
func NewPartItem(
	id uuid.UUID,
	name string,
	description string,
	price float64,
	stockQuantity int64,
	category Category,
	dimensions *Dimensions,
	manufacturer *Manufacturer,
	tags []string,
	metadata map[string]*Value,
	createdAt *time.Time,
	updatedAt *time.Time,
) *PartItem {
	return &model.PartItem{
		ID:            id,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		Metadata:      metadata,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
