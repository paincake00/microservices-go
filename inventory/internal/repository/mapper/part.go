package mapper

import (
	"encoding/json"
	"errors"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func ToModel(p *entity.Part) (*model.Part, error) {
	metadata, err := metadataToBytes(p.Metadata)
	if err != nil {
		return nil, err
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
		Metadata:  metadata,
		CreatedAt: *p.CreatedAt,
		UpdatedAt: *p.UpdatedAt,
	}, nil
}

func FromModel(p *model.Part) (*entity.Part, error) {
	metadata, err := metadataFromBytes(p.Metadata)
	if err != nil {
		return nil, err
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
		Metadata:  metadata,
		CreatedAt: &p.CreatedAt,
		UpdatedAt: &p.UpdatedAt,
	}, nil
}

func metadataToBytes(metadata map[string]*entity.Value) ([]byte, error) {
	b, err := json.Marshal(metadata)
	if err != nil {
		return nil, model.ErrSerializeMetadata
	}
	return b, nil
}

func metadataFromBytes(jsonb []byte) (map[string]*entity.Value, error) {
	var metadata map[string]*entity.Value
	err := json.Unmarshal(jsonb, &metadata)
	if err != nil {
		if errors.Is(err, entity.ErrValueKind) {
			return nil, err
		}
		return nil, model.ErrUnserializeMetadata
	}
	return metadata, nil
}
