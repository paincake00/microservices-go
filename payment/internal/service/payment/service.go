package payment

import "github.com/google/uuid"

type UuidGenerator interface {
	NewV7() (uuid.UUID, error)
}

type Service struct {
	uuidGenerator UuidGenerator
}

func NewService(generator UuidGenerator) *Service {
	return &Service{
		uuidGenerator: generator,
	}
}
