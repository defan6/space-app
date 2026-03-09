package inventory

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/client/converter"
	inventorymodel "github.com/defan6/space-app/order-service/internal/client/inventory/model"
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func (ic *inventoryClient) ExternalGetPart(ctx context.Context, partUUID uuid.UUID) (*inventorymodel.PartItem, error) {
	op := "inventory-client#ExternalGetPart"

	partReq := &inventoryV1.GetPartRequest{
		Uuid: partUUID.String(),
	}

	partExtRes, err := ic.GetPart(ctx, partReq)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	partRes, err := converter.FromInventoryExternalGetPartResponse(partExtRes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return partRes, nil
}
