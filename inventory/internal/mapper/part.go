package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func ToProto(p *entity.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      CategoryToProto(p.Category),
		Dimensions: &inventoryv1.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  parseMetadata(p.Metadata),
		CreatedAt: timestamppb.New(*p.CreatedAt),
		UpdatedAt: timestamppb.New(*p.UpdatedAt),
	}
}

func CategoryToProto(category entity.Category) inventoryv1.Category {
	switch category {
	case entity.CATEGORY_ENGINE:
		return inventoryv1.Category_CATEGORY_ENGINE
	case entity.CATEGORY_FUEL:
		return inventoryv1.Category_CATEGORY_FUEL
	case entity.CATEGORY_PORTHOLE:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case entity.CATEGORY_WING:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

func CategoryFromProto(category inventoryv1.Category) entity.Category {
	switch category {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return entity.CATEGORY_ENGINE
	case inventoryv1.Category_CATEGORY_FUEL:
		return entity.CATEGORY_FUEL
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return entity.CATEGORY_PORTHOLE
	case inventoryv1.Category_CATEGORY_WING:
		return entity.CATEGORY_WING
	default:
		return entity.CATEGORY_UNSPECIFIED
	}
}

func parseMetadata(m map[string]*entity.Value) map[string]*inventoryv1.Value {
	res := make(map[string]*inventoryv1.Value)

	for k, v := range m {
		var value inventoryv1.Value
		switch v.GetKind() {
		case entity.StringValue:
			value = inventoryv1.Value{
				Kind: &inventoryv1.Value_StringValue{
					StringValue: v.GetStringValue(),
				},
			}
		case entity.Int64Value:
			value = inventoryv1.Value{
				Kind: &inventoryv1.Value_Int64Value{
					Int64Value: v.GetInt64Value(),
				},
			}
		case entity.Float64Value:
			value = inventoryv1.Value{
				Kind: &inventoryv1.Value_DoubleValue{
					DoubleValue: v.GetFloat64Value(),
				},
			}
		case entity.BoolValue:
			value = inventoryv1.Value{
				Kind: &inventoryv1.Value_BoolValue{
					BoolValue: v.GetBoolValue(),
				},
			}
		}
		res[k] = &value
	}

	return res
}
