package part

import (
	"context"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

func (p *Service) GetPart(ctx context.Context, id string) (*entity.Part, error) {
	part, err := p.partRepo.GetPart(ctx, id)
	if err != nil {
		return nil, err
	}
	return part, nil
}
