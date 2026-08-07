package part

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func (p *Repository) GetPart(id string) (*entity.Part, error) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	part, ok := p.parts[id]
	if !ok {
		return nil, model.ErrNotFound
	}

	res, err := mapper.FromModel(part)
	if err != nil {
		return nil, err
	}
	return res, nil
}
