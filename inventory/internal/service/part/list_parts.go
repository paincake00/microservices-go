package part

import "github.com/paincake00/microservices-go/inventory/internal/entity"

func (p *Service) ListParts(filter entity.ListPartsFilter) ([]*entity.Part, error) {
	parts, err := p.partRepo.ListParts(filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
