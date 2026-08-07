package part

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
)

func (p *Repository) Create(part *entity.Part) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	partModel, err := mapper.ToModel(part)
	if err != nil {
		return err
	}

	p.parts[part.Uuid] = partModel

	return nil
}
