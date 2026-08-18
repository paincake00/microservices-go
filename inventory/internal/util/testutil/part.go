package testutil

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
)

func NewPart(id string, category enum.Category, metadata map[string]*entity.Value) *entity.Part {
	createdAt := gofakeit.Date().Truncate(time.Millisecond)
	updatedAt := gofakeit.DateRange(createdAt, createdAt.AddDate(0, 0, 1)).Truncate(time.Millisecond)

	part := &entity.Part{
		Uuid:          id,
		Name:          gofakeit.Name(),
		Description:   gofakeit.ProductDescription(),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: gofakeit.Int64(),
		Category:      category,
		Dimensions: &entity.Dimensions{
			Height: gofakeit.Float64(),
			Width:  gofakeit.Float64(),
			Length: gofakeit.Float64(),
			Weight: gofakeit.Float64(),
		},
		Manufacturer: &entity.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{
			gofakeit.Word(),
			gofakeit.Word(),
			gofakeit.Word(),
		},
		Metadata:  metadata,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	return part
}

func AllCategories() ([]enum.Category, int) {
	res := make([]enum.Category, 0, len(enum.CategoryAll))

	for _, name := range enum.CategoryAll {
		if enum.Category(name) == enum.CategoryUnspecified {
			continue
		}

		res = append(res, enum.Category(name))
	}

	return res, len(res)
}

func RandomMetadata() map[string]*entity.Value {
	return map[string]*entity.Value{
		"color":   entity.NewStringValue(gofakeit.Color()),
		"count":   entity.NewInt64Value(gofakeit.Int64()),
		"rating":  entity.NewFloat64Value(gofakeit.Float64()),
		"enabled": entity.NewBoolValue(gofakeit.Bool()),
	}
}
