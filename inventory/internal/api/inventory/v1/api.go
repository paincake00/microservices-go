package v1

import (
	"github.com/paincake00/microservices-go/inventory/internal/service"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer

	partService service.IPartService
}

func NewInventoryHandler(partService service.IPartService) *InventoryHandler {
	return &InventoryHandler{
		partService: partService,
	}
}
