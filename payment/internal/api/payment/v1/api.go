package v1

import (
	"github.com/paincake00/microservices-go/payment/internal/service"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

type PaymentHandler struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.IPaymentService
}

func NewPaymentHandler(paymentService service.IPaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}
