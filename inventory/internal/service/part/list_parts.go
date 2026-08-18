package part

import (
	"context"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

func (p *Service) ListParts(ctx context.Context, filter entity.ListPartsFilter) ([]*entity.Part, error) {
	parts, err := p.partRepo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
