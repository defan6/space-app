package inventory

import (
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type inventoryClient struct {
	inventoryV1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *inventoryClient {
	return &inventoryClient{
		InventoryServiceClient: inventoryV1.NewInventoryServiceClient(conn),
	}
}
