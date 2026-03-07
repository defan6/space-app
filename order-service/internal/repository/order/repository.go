package order

import (
	"sync"

	repomodel "github.com/defan6/space-app/order-service/internal/repository/model"
	"github.com/google/uuid"
)

type inMemRepo struct {
	cache map[uuid.UUID]*repomodel.Order
	mu    *sync.RWMutex
}

func newInMemRepo() *inMemRepo {
	return &inMemRepo{
		cache: make(map[uuid.UUID]*repomodel.Order),
		mu:    &sync.RWMutex{},
	}
}
