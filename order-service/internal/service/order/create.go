package serviceorder

import (
	"context"
	"fmt"
	"sync"

	"github.com/defan6/space-app/order-service/internal/model"
	repoconverter "github.com/defan6/space-app/order-service/internal/repository/converter"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func (dos *defaultOrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	op := "order-service#CreateOrder"

	checkInvRes, err := dos.checkInventory(ctx, req.PartItems)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	order := repoconverter.ConvertFromCreateOrderRequestToRepoOrder(req)
	order.PartsItems = repoconverter.ConvertFromPartInfosToRepoParts(checkInvRes.PartItems)
	order.TotalPrice = checkInvRes.TotalPrice

	createdOrder, err := dos.repo.Save(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return repoconverter.ConvertFromOrderToCreateOrderResponse(createdOrder), nil
}

func (dos *defaultOrderService) checkInventory(ctx context.Context, parts []model.PartItemRequest) (*model.CheckInventoryResponse, error) {
	op := "order-service#CheckInventory"

	partsMap := make(map[uuid.UUID]model.PartItemRequest)
	for _, p := range parts {
		partsMap[p.PartUUID] = p
	}

	fetchedPartsMap, err := dos.fetchParts(parts, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	partsWithPrice := make([]model.PartInfo, 0, len(fetchedPartsMap))
	var totalPrice float64

	for id, partInfo := range fetchedPartsMap {
		requestedQty := partsMap[id].Quantity
		totalPrice += float64(requestedQty) * partInfo.Price

		if partInfo.Quantity < requestedQty {
			return nil, fmt.Errorf("%s: insufficient stock for part %s", op, id.String())
		}

		partWithPrice := model.PartInfo{
			PartUUID: id,
			Quantity: requestedQty,
			Price:    partInfo.Price,
		}
		partsWithPrice = append(partsWithPrice, partWithPrice)
	}

	return &model.CheckInventoryResponse{
		TotalPrice: totalPrice,
		PartItems:  partsWithPrice,
	}, nil
}

func (dos *defaultOrderService) fetchParts(parts []model.PartItemRequest, ctx context.Context) (map[uuid.UUID]*model.PartInfo, error) {
	op := "order-service#fetchParts"

	fetchedPartsMap := make(map[uuid.UUID]*model.PartInfo)
	mu := &sync.Mutex{}

	egroup, ectx := errgroup.WithContext(ctx)
	egroup.SetLimit(5)

	for _, partItem := range parts {
		egroup.Go(func() error {
			partFromInventory, err := dos.inventoryClient.ExternalGetPart(ectx, partItem.PartUUID)
			if err != nil {
				return fmt.Errorf("failed to fetch part %s: %w", partItem.PartUUID.String(), err)
			}

			mu.Lock()
			fetchedPartsMap[partItem.PartUUID] = &model.PartInfo{
				PartUUID: partFromInventory.ID,
				Quantity: partFromInventory.StockQuantity,
				Price:    partFromInventory.Price,
			}
			mu.Unlock()

			return nil
		})
	}

	if err := egroup.Wait(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return fetchedPartsMap, nil
}
