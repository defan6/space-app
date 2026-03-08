package repoorder

import (
	"context"
	"log"
	"time"

	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

func (imr *inMemRepo) Save(_ context.Context, order *repomodel.Order) (*repomodel.Order, error) {
	op := "order-repo#CreateOrder"
	ouuid := uuid.New()
	order.OrderUUID = ouuid
	order.CreatedAt = time.Now()
	imr.mu.Lock()
	imr.cache[ouuid] = order
	imr.mu.Unlock()
	log.Printf("%s: %#v", op, order)
	return order, nil
}
