package order

import (
	"context"
	"log"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

func (or *Service) Pay(
	ctx context.Context,
	orderUUID string,
	payMethod enum.PaymentMethod,
) (entity.PaymentTransaction, error) {
	order, err := or.orderStorage.Get(ctx, orderUUID)
	if err != nil {
		return entity.PaymentTransaction{}, err
	}

	payResponse, err := or.paymentService.PayOrder(ctx, orderUUID, order.UserUuid, payMethod)
	if err != nil {
		return entity.PaymentTransaction{}, err
	}

	// Меняем поля у копии
	order.TransactionUuid = &payResponse.TransactionUuid
	order.PaymentMethod = &payMethod
	order.Status = enum.Paid

	// Обновляем в карте заказ с новыми полями
	errUpdate := or.orderStorage.Update(ctx, order)
	if errUpdate != nil {
		return entity.PaymentTransaction{}, errUpdate
	}

	log.Printf("Pay order: %s", order.OrderUuid)

	return entity.PaymentTransaction{TransactionUuid: *order.TransactionUuid}, nil
}
