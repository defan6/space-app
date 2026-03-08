package order

import (
	"context"

	"github.com/google/uuid"
)

func (imr *inMemRepo) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	_, ok := imr.cache[uuid]
	return ok, nil
}
