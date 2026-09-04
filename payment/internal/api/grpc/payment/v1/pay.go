package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

func (p *PaymentHandler) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (
	*paymentv1.PayOrderResponse,
	error,
) {
	id, err := p.paymentService.PayOrder()
	if err != nil {
		log.Printf("Транзакция отменена. transaction_uuid не создан: %v", err)

		return nil, status.Errorf(codes.Internal, "Транзакция отменена. transaction_uuid не создан")
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: id,
	}, nil
}
