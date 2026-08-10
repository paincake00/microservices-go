package mapper

import (
	"github.com/paincake00/microservices-go/order/internal/entity"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func FromProto(part *inventoryv1.Part) *entity.Part {
	return &entity.Part{
		Price: part.Price,
	}
}
