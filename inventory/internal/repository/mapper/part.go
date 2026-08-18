package mapper

import (
	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func ToModel(p *entity.Part) (*model.Part, error) {
	if p == nil {
		return nil, model.ErrNilPart
	}

	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      p.Category,
		Dimensions: entity.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: entity.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  p.Metadata,
		CreatedAt: *p.CreatedAt,
		UpdatedAt: *p.UpdatedAt,
	}, nil
}

func FromModel(p *model.Part) (*entity.Part, error) {
	if p == nil {
		return nil, model.ErrNilPart
	}

	return &entity.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      p.Category,
		Dimensions: &entity.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: &entity.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  p.Metadata,
		CreatedAt: &p.CreatedAt,
		UpdatedAt: &p.UpdatedAt,
	}, nil
}
