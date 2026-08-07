package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paincake00/microservices-go/inventory/internal/mapper"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func (i *InventoryHandler) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (
	*inventoryv1.GetPartResponse,
	error,
) {
	id := req.GetUuid()

	p, err := i.partService.GetPart(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Printf("Part %s not found", id)

			return nil, status.Errorf(codes.NotFound, "not found")
		}

		log.Printf("Error getting part %s: %s", id, err)

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &inventoryv1.GetPartResponse{
		Part: mapper.ToProto(p),
	}, nil
}
