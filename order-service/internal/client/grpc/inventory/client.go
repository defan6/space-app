package inventory

import (
	inventoryV1 "github.com/defan6/space-app/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type inventoryClient struct {
	client inventoryV1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *inventoryClient {
	client := inventoryV1.NewInventoryServiceClient(conn)
	return &inventoryClient{
		client: client,
	}
}
