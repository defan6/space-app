package serviceorder

import (
	"github.com/defan6/space-app/order-service/internal/client"
	"github.com/defan6/space-app/order-service/internal/repository"
)

type defaultOrderService struct {
	repo            repository.Repository
	paymentClient   client.PaymentClient
	inventoryClient client.InventoryClient
}

func NewService(repo repository.Repository, pclient client.PaymentClient, iclient client.InventoryClient) *defaultOrderService {
	return &defaultOrderService{
		repo:            repo,
		paymentClient:   pclient,
		inventoryClient: iclient,
	}
}
