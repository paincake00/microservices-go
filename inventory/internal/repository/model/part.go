package model

import (
	"time"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
)

type Part struct {
	Uuid          string                   `bson:"_id,omitempty"`
	Name          string                   `bson:"name"`
	Description   string                   `bson:"description"`
	Price         float64                  `bson:"price"`
	StockQuantity int64                    `bson:"stock_quantity"`
	Category      enum.Category            `bson:"category"`
	Dimensions    entity.Dimensions        `bson:"dimensions"`
	Manufacturer  entity.Manufacturer      `bson:"manufacturer"`
	Tags          []string                 `bson:"tags"`
	Metadata      map[string]*entity.Value `bson:"metadata,omitempty"`
	CreatedAt     time.Time                `bson:"created_at"`
	UpdatedAt     time.Time                `bson:"updated_at,omitempty"`
}
