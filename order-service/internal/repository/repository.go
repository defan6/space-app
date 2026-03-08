package repository

import (
	"context"

	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, order *repomodel.Order) (*repomodel.Order, error)
	Get(ctx context.Context, uuid uuid.UUID) (*repomodel.Order, error)
	GetOrders(ctx context.Context) ([]*repomodel.Order, error)
	Update(ctx context.Context, order *repomodel.Order) (*repomodel.Order, error)
	Exists(ctx context.Context, uuid uuid.UUID) (bool, error)
}
