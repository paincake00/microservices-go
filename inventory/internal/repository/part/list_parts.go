package part

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
)

func (p *Repository) ListParts(filter entity.ListPartsFilter) ([]*entity.Part, error) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	res := make([]*entity.Part, 0)

	emptyFilter := filter.IsEmpty()

	for _, pp := range p.parts {
		part, err := mapper.FromModel(pp)
		if err != nil {
			return nil, err
		}
		if emptyFilter || filter.Match(part) {
			res = append(res, part)
		}
	}

	return res, nil
}
