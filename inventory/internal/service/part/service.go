package part

import "github.com/paincake00/microservices-go/inventory/internal/repository"

type Service struct {
	partRepo repository.IPartRepository
}

func NewService(partRepo repository.IPartRepository) *Service {
	return &Service{
		partRepo: partRepo,
	}
}
