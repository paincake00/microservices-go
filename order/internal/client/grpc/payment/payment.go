package payment

import (
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

type Service struct {
	grpcPaymentClient paymentv1.PaymentServiceClient
}

func NewService(grpcPaymentClient paymentv1.PaymentServiceClient) *Service {
	return &Service{
		grpcPaymentClient: grpcPaymentClient,
	}
}
