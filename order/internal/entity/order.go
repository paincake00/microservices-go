package entity

import "github.com/paincake00/microservices-go/order/internal/entity/enum"

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid *string
	PaymentMethod   *enum.PaymentMethod
	Status          enum.OrderStatus
}
