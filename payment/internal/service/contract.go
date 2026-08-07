package service

type IPaymentService interface {
	PayOrder() (string, error)
}
