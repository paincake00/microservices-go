package inventory

import (
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

type Service struct {
	grpcInventoryClient inventoryv1.InventoryServiceClient
}

func NewService(grpcInventoryClient inventoryv1.InventoryServiceClient) *Service {
	return &Service{
		grpcInventoryClient: grpcInventoryClient,
	}
}
