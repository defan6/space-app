package serviceorder

import (
	"context"
	"fmt"
	"sync"

	"github.com/defan6/space-app/order-service/internal/model"
	"github.com/defan6/space-app/order-service/internal/repository/converter"
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func (dos *defaultOrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	op := "order-service#CreateOrder"
	checkInvRes, err := dos.checkInventory(ctx, req.PartItems)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	order := converter.ConvertFromCreateOrderRequestToRepoOrder(req)
	order.PartsItems = converter.ConvertFromPartsToRepoParts(checkInvRes.PartItems)
	order.TotalPrice = checkInvRes.TotalPrice
	createdOrder, err := dos.repo.Save(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return converter.ConvertFromOrderToCreateOrderResponse(createdOrder), nil
}

func (dos *defaultOrderService) checkInventory(ctx context.Context, parts []model.PartItem) (*model.CheckInventoryResponse, error) {
	op := "order-service#CheckInventory"
	partsMap := make(map[uuid.UUID]model.PartItem)

	for _, p := range parts {
		partsMap[p.PartUUID] = p
	}
	fetchedPartsMap, err := dos.fetchParts(parts, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	partsWithPrice := make([]model.PartItem, 0, len(fetchedPartsMap))
	var totalPrice float64
	for id, part := range fetchedPartsMap {
		totalPrice += float64(partsMap[id].Quantity) * fetchedPartsMap[id].Price
		if part.Quantity < partsMap[id].Quantity {
			return nil, fmt.Errorf("%s (partID: %s):%w", op, id.String(), err)
		}
		partWithPrice := model.PartItem{
			PartUUID: id,
			Quantity: partsMap[id].Quantity,
			Price:    fetchedPartsMap[id].Price,
		}
		partsWithPrice = append(partsWithPrice, partWithPrice)
	}

	return &model.CheckInventoryResponse{
		TotalPrice: totalPrice,
		PartItems:  partsWithPrice,
	}, nil
}

func (dos *defaultOrderService) fetchParts(parts []model.PartItem, ctx context.Context) (map[uuid.UUID]model.PartItem, error) {
	op := "order-service#fetchParts"
	fetchedPartsMap := make(map[uuid.UUID]model.PartItem)
	mu := &sync.Mutex{}

	egroup, ectx := errgroup.WithContext(ctx)
	egroup.SetLimit(5)
	for _, partItem := range parts {
		egroup.Go(func() error {
			invReq := &inventoryV1.GetPartRequest{
				Uuid: partItem.PartUUID.String(),
			}
			invRes, err := dos.inventoryClient.GetPart(ectx, invReq)
			if err != nil {
				return fmt.Errorf("failed to fetch part %s:%w", partItem.PartUUID.String(), err)
			}
			mu.Lock()
			fetchedPartsMap[partItem.PartUUID] = model.PartItem{
				PartUUID: partItem.PartUUID,
				Quantity: invRes.StockQuantity,
				Price:    invRes.Price,
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
