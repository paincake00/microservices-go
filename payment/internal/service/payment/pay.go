package payment

import (
	"log"
)

func (p *Service) PayOrder() (string, error) {
	id, err := p.uuidGenerator.NewV7()
	if err != nil {
		return "", err
	}
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", id.String())

	return id.String(), nil
}
