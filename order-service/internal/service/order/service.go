package order

import (
	"github.com/defan6/space-app/order-service/internal/repository"
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/defan6/space-app/shared/pkg/proto/payment/v1"
)

type defaultOrderService struct {
	repo            repository.Repository
	paymentClient   paymentV1.PaymentServiceClient
	inventoryClient inventoryV1.InventoryServiceClient
}

func NewService(repo repository.Repository, pclient paymentV1.PaymentServiceClient, iclient inventoryV1.InventoryServiceClient) *defaultOrderService {
	return &defaultOrderService{
		repo:            repo,
		paymentClient:   pclient,
		inventoryClient: iclient,
	}
}
