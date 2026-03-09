package paymentmodel

import (
	"github.com/google/uuid"
)

// PayOrderInfo — информация об оплате заказа.
type PayOrderInfo struct {
	TransactionUUID uuid.UUID
}
