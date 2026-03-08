package order

import "github.com/defan6/space-app/order-service/internal/service"

type orderAPI struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *orderAPI {
	return &orderAPI{
		service: service,
	}
}
