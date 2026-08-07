package part

import "github.com/paincake00/microservices-go/inventory/internal/entity"

func (p *Service) GetPart(id string) (*entity.Part, error) {
	part, err := p.partRepo.GetPart(id)
	if err != nil {
		return nil, err
	}
	return part, nil
}
