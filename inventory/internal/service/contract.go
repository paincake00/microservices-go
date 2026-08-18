package service

import (
	"context"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

type IPartService interface {
	GetPart(ctx context.Context, partUUID string) (*entity.Part, error)
	ListParts(ctx context.Context, filter entity.ListPartsFilter) ([]*entity.Part, error)
}
