package part

import (
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/util/uuidutil"
)

// InitParts - создание исходных деталей при старте сервиса
// numParts - число создаваемых в начале деталей
func (p *Service) InitParts(numParts int) error {
	categories, length := allCategories()

	for range numParts {
		id, err := uuidutil.GetNewUUID()
		if err != nil {
			log.Fatal(err)
		}

		//nolint:gosec // random is used only for test data generation
		category := categories[rand.Intn(length)]

		createdAt := gofakeit.Date()
		updatedAt := gofakeit.DateRange(createdAt, createdAt.AddDate(0, 0, 1))

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
			Metadata:  randomMetadata(),
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}

		log.Printf("Add new part with UUID: %s", part.Uuid)

		err = p.partRepo.Create(part)
		if err != nil {
			return err
		}
	}

	return nil
}

func allCategories() ([]entity.Category, int) {
	res := make([]entity.Category, 0, len(entity.CategoryName))

	for name := range entity.CategoryName {
		if entity.Category(name) == entity.CATEGORY_UNSPECIFIED {
			continue
		}

		res = append(res, entity.Category(name))
	}

	return res, len(res)
}

func randomMetadata() map[string]*entity.Value {
	return map[string]*entity.Value{
		"color":   entity.NewStringValue(gofakeit.Color()),
		"count":   entity.NewInt64Value(gofakeit.Int64()),
		"rating":  entity.NewFloat64Value(gofakeit.Float64()),
		"enabled": entity.NewBoolValue(gofakeit.Bool()),
	}
}
