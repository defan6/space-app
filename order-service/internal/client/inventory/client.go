package inventory

import (
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type inventoryClient struct {
	inventoryV1.InventoryServiceClient
}

// NewInventoryClient создаёт нового gRPC клиента для подключения к inventory service.
func NewInventoryClient(conn *grpc.ClientConn) *inventoryClient {
	return &inventoryClient{
		InventoryServiceClient: inventoryV1.NewInventoryServiceClient(conn),
	}
}
