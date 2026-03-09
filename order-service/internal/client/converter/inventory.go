package converter

import (
	"fmt"
	"time"

	inventorymodel "github.com/defan6/space-app/order-service/internal/client/inventory/model"
	paymentmodel "github.com/defan6/space-app/order-service/internal/client/payment/model"
	"github.com/defan6/space-app/order-service/internal/model"
	orderV1 "github.com/defan6/space-app/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// PartItem — основной конвертер
// =============================================================================

// FromInventoryExternalGetPartResponse конвертирует gRPC GetPartResponse → inventorymodel.PartItem.
func FromInventoryExternalGetPartResponse(resp *inventoryV1.GetPartResponse) (*inventorymodel.PartItem, error) {
	op := "inventory-converter#FromInventoryExternalGetPartResponse"

	id, err := uuid.Parse(resp.Uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse UUID: %w", op, err)
	}

	return &inventorymodel.PartItem{
		ID:            id,
		Name:          resp.Name,
		Description:   resp.Description,
		Price:         resp.Price,
		StockQuantity: resp.StockQuantity,
		Category:      FromInventoryCategory(resp.Category),
		Dimensions:    FromInventoryDimensions(resp.Dimensions),
		Manufacturer:  FromInventoryManufacturer(resp.Manufacturer),
		Tags:          FromInventoryTags(resp.Tags),
		Metadata:      FromInventoryMetadata(resp.Metadata),
		CreatedAt:     FromInventoryTimestamp(resp.CreatedAt),
		UpdatedAt:     FromInventoryTimestamp(resp.UpdatedAt),
	}, nil
}

// =============================================================================
// Category
// =============================================================================

// FromInventoryCategory конвертирует gRPC Category enum → model.Category.
func FromInventoryCategory(cat inventoryV1.Category) inventorymodel.Category {
	switch cat {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return inventorymodel.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return inventorymodel.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return inventorymodel.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return inventorymodel.CategoryWing
	default:
		return inventorymodel.CategoryUnknown
	}
}

// ToInventoryCategory конвертирует model.Category → gRPC Category enum.
func ToInventoryCategory(cat inventorymodel.Category) inventoryV1.Category {
	switch cat {
	case inventorymodel.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case inventorymodel.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case inventorymodel.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case inventorymodel.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	}
}

// =============================================================================
// Dimensions
// =============================================================================

// FromInventoryDimensions конвертирует gRPC Dimensions → model.Dimensions.
func FromInventoryDimensions(dim *inventoryV1.Dimensions) *inventorymodel.Dimensions {
	if dim == nil {
		return nil
	}

	return &inventorymodel.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// ToInventoryDimensions конвертирует model.Dimensions → gRPC Dimensions.
func ToInventoryDimensions(dim *inventorymodel.Dimensions) *inventoryV1.Dimensions {
	if dim == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// =============================================================================
// Manufacturer
// =============================================================================

// FromInventoryManufacturer конвертирует gRPC Manufacturer → model.Manufacturer.
func FromInventoryManufacturer(mfg *inventoryV1.Manufacturer) *inventorymodel.Manufacturer {
	if mfg == nil {
		return nil
	}

	return &inventorymodel.Manufacturer{
		Name:    mfg.Name,
		Country: mfg.Country,
		Website: mfg.Website,
	}
}

// ToInventoryManufacturer конвертирует model.Manufacturer → gRPC Manufacturer.
func ToInventoryManufacturer(mfg *inventorymodel.Manufacturer) *inventoryV1.Manufacturer {
	if mfg == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    mfg.Name,
		Country: mfg.Country,
		Website: mfg.Website,
	}
}

// =============================================================================
// Tags
// =============================================================================

// FromInventoryTags конвертирует gRPC []string → model.Tags.
func FromInventoryTags(tags []string) []string {
	if tags == nil {
		return nil
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// ToInventoryTags конвертирует model.Tags → gRPC []string.
func ToInventoryTags(tags []string) []string {
	if tags == nil {
		return nil
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// =============================================================================
// Metadata
// =============================================================================

// FromInventoryMetadata конвертирует gRPC map[string]*Value → model.Metadata.
func FromInventoryMetadata(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*model.Value, len(metadata))
	for key, val := range metadata {
		if val == nil {
			continue
		}
		result[key] = FromInventoryValue(val)
	}
	return result
}

// ToInventoryMetadata конвертирует model.Metadata → gRPC map[string]*Value.
func ToInventoryMetadata(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, val := range metadata {
		if val == nil {
			continue
		}
		result[key] = ToInventoryValue(val)
	}
	return result
}

// FromInventoryValue конвертирует gRPC Value → model.Value.
func FromInventoryValue(val *inventoryV1.Value) *model.Value {
	if val == nil {
		return nil
	}

	return &model.Value{
		StringValue: val.StringValue,
		Int64Value:  val.Int64Value,
		DoubleValue: val.DoubleValue,
		BoolValue:   val.BoolValue,
	}
}

// ToInventoryValue конвертирует model.Value → gRPC Value.
func ToInventoryValue(val *model.Value) *inventoryV1.Value {
	if val == nil {
		return nil
	}

	return &inventoryV1.Value{
		StringValue: val.StringValue,
		Int64Value:  val.Int64Value,
		DoubleValue: val.DoubleValue,
		BoolValue:   val.BoolValue,
	}
}

// =============================================================================
// Timestamp
// =============================================================================

// FromInventoryTimestamp конвертирует gRPC Timestamp → *time.Time.
func FromInventoryTimestamp(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}

// ToInventoryTimestamp конвертирует *time.Time → gRPC Timestamp.
func ToInventoryTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

// =============================================================================
// OrderStatus — общие конвертеры (используются в API layer)
// =============================================================================

// ConvertFromAPIStatusToStatus конвертирует API OrderStatus → model.OrderStatus.
func ConvertFromAPIStatusToStatus(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCancelled
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPaid
	case orderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPendingPayment
	default:
		return model.OrderStatusPendingPayment
	}
}

// ConvertFromStatusToAPIStatus конвертирует model.OrderStatus → API OrderStatus.
func ConvertFromStatusToAPIStatus(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

// =============================================================================
// PaymentMethod — общие конвертеры (используются в API layer)
// =============================================================================

// ConvertFromAPIPaymentMethodToPaymentMethod конвертирует API PaymentMethod → model.PaymentMethod.
func ConvertFromAPIPaymentMethodToPaymentMethod(pm orderV1.PaymentMethod) (model.PaymentMethod, error) {
	switch pm {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCredit, nil
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnknown, nil
	}
}

// ConvertFromPaymentMethodToAPIPaymentMethod конвертирует model.PaymentMethod → API PaymentMethod.
func ConvertFromPaymentMethodToAPIPaymentMethod(pm model.PaymentMethod) orderV1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodCredit:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

// =============================================================================
// Payment — конвертеры для payment service
// =============================================================================

// FromPaymentExternalPayOrderResponse конвертирует gRPC PayOrderResponse → paymentmodel.PayOrderInfo.
func FromPaymentExternalPayOrderResponse(resp *paymentV1.PayOrderResponse) (*paymentmodel.PayOrderInfo, error) {
	op := "payment-converter#FromPaymentExternalPayOrderResponse"

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse transaction UUID: %w", op, err)
	}

	return &paymentmodel.PayOrderInfo{
		TransactionUUID: transactionUUID,
	}, nil
}

// FromPaymentMethod конвертирует gRPC PaymentMethod → model.PaymentMethod.
func FromPaymentMethod(pm paymentV1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCredit
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
