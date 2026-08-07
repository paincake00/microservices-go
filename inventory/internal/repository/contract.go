package repository

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

type IPartRepository interface {
	GetPart(string) (*entity.Part, error)
	ListParts(filter entity.ListPartsFilter) ([]*entity.Part, error)
	Create(part *entity.Part) error
}
