package repoorder

import (
	"context"
	"fmt"
	"time"

	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/samber/lo"
)

func (imr *inMemRepo) Update(ctx context.Context, order *repomodel.Order) (*repomodel.Order, error) {
	op := "order-repo#Update"

	existingOrder, err := imr.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	imr.update(ctx, existingOrder, order)

	imr.mu.Lock()
	imr.cache[existingOrder.OrderUUID] = existingOrder
	defer imr.mu.Unlock()

	return existingOrder, nil

}

func (imr *inMemRepo) update(ctx context.Context, existingOrder *repomodel.Order, updatedOrder *repomodel.Order) {
	if updatedOrder.UserUUID.String() != "" {
		existingOrder.UserUUID = updatedOrder.UserUUID
	}
	if updatedOrder.TransactionUUID.String() != "" {
		existingOrder.TransactionUUID = updatedOrder.TransactionUUID
	}
	if updatedOrder.TotalPrice != 0 {
		existingOrder.TotalPrice = updatedOrder.TotalPrice
	}
	if updatedOrder.Status != "" {
		existingOrder.Status = updatedOrder.Status
	}
	if updatedOrder.PaymentMethod != "" {
		existingOrder.PaymentMethod = updatedOrder.PaymentMethod
	}
	if updatedOrder.PartsItems != nil {
		existingOrder.PartsItems = updatedOrder.PartsItems
	}
	existingOrder.UpdatedAt = lo.ToPtr(time.Now())
}
