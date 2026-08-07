package payment

import (
	"log"

	"github.com/google/uuid"
)

func (p *Service) PayOrder() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", id.String())

	return id.String(), nil
}
