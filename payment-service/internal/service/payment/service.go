package payment

type defaultPaymentService struct{}

func NewDefaultPaymentService() *defaultPaymentService {
	return &defaultPaymentService{}
}
