package client

import (
	"context"

	inventorymodel "github.com/defan6/space-app/order-service/internal/client/inventory/model"
	paymentmodel "github.com/defan6/space-app/order-service/internal/client/payment/model"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

type InventoryClient interface {
	ExternalGetPart(ctx context.Context, partUUID uuid.UUID) (*inventorymodel.PartItem, error)
}

type PaymentClient interface {
	ExternalPayOrder(ctx context.Context, orderUUID uuid.UUID, req *paymentV1.PayOrderRequest) (*paymentmodel.PayOrderInfo, error)
}
