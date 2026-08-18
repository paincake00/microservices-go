package repository

import (
	"context"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

type IPartRepository interface {
	GetPart(ctx context.Context, partUUID string) (*entity.Part, error)
	ListParts(ctx context.Context, filter entity.ListPartsFilter) ([]*entity.Part, error)
	Create(ctx context.Context, part *entity.Part) (string, error)
	DeleteAll(ctx context.Context) error
}
