package order

import (
	"context"
	"fmt"

	"github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

func (imr *inMemRepo) GetOrder(ctx context.Context, uuid uuid.UUID) (*repomodel.Order, error) {
	op := "order-repo#GetOrder"
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	order, ok := imr.cache[uuid]
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, repomodel.ErrNotFound)
	}
	return order, nil
}
