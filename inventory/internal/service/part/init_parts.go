package part

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
	"github.com/paincake00/microservices-go/inventory/internal/util/uuidutil"
)

// InitParts - создание исходных деталей при старте сервиса
// numParts - число создаваемых в начале деталей
func (p *Service) InitParts(ctx context.Context, numParts int) error {
	// Очистка предыдущих записей из коллекции
	err := p.partRepo.DeleteAll(ctx)
	if err != nil {
		return err
	}

	// Создание фейковых категорий
	categories, length := allCategories()

	// Создание фейковых деталей
	for range numParts {
		id, err := uuidutil.GetNewUUID()
		if err != nil {
			log.Fatal(err)
		}

		//nolint:gosec // random is used only for test data generation
		category := categories[rand.Intn(length)]

		part := createFakedPart(id, category)

		id, err = p.partRepo.Create(ctx, part)
		if err != nil {
			return err
		}
		log.Printf("Add new part with UUID: %s", id)
	}

	return nil
}

func createFakedPart(id string, category enum.Category) *entity.Part {
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
		Metadata:  randomMetadata(),
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	return part
}

func allCategories() ([]enum.Category, int) {
	res := make([]enum.Category, 0, len(enum.CategoryAll))

	for _, name := range enum.CategoryAll {
		if enum.Category(name) == enum.CategoryUnspecified {
			continue
		}

		res = append(res, enum.Category(name))
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
