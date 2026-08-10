package model

import (
	"time"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
)

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      enum.Category
	Dimensions    entity.Dimensions
	Manufacturer  entity.Manufacturer
	Tags          []string
	Metadata      []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
