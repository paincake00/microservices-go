package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/paincake00/microservices-go/inventory/internal/api/grpc/inventory/v1"
	"github.com/paincake00/microservices-go/inventory/internal/service"
	"github.com/paincake00/microservices-go/platform/pkg/grpc/health"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func NewRouter(srv *grpc.Server, partService service.IPartService) {
	inventoryHandler := v1.NewInventoryHandler(partService)

	inventoryv1.RegisterInventoryServiceServer(srv, inventoryHandler)

	reflection.Register(srv)

	// Регистрация эндпоинтов для healthcheck
	health.RegisterServer(srv)
}
