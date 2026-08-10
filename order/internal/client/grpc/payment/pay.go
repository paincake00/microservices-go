package payment

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	"github.com/paincake00/microservices-go/order/internal/mapper"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

// PayOrder returns a transactionUUID or error
func (s *Service) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod enum.PaymentMethod) (
	entity.PaymentTransaction,
	error,
) {
	resp, err := s.grpcPaymentClient.PayOrder(
		ctx, &paymentv1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: mapper.PaymentMethodToPaymentProto(paymentMethod),
		},
	)
	if err != nil {
		return entity.PaymentTransaction{}, err
	}
	return entity.PaymentTransaction{
		TransactionUuid: resp.TransactionUuid,
	}, nil
}
