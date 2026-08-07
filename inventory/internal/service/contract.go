package service

import "github.com/paincake00/microservices-go/inventory/internal/entity"

type IPartService interface {
	GetPart(string) (*entity.Part, error)
	ListParts(filter entity.ListPartsFilter) ([]*entity.Part, error)
}
