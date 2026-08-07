package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/mapper"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func (i *InventoryHandler) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	protoFilter := req.GetFilter()

	categories := make([]entity.Category, 0)

	for _, c := range protoFilter.GetCategories() {
		categories = append(categories, mapper.CategoryFromProto(c))
	}

	filter := entity.NewListPartsFilter(
		protoFilter.GetUuids(), protoFilter.GetNames(), categories,
		protoFilter.GetManufacturerCountries(), protoFilter.GetTags(),
	)

	res := make([]*inventoryv1.Part, 0)
	parts, err := i.partService.ListParts(filter)
	if err != nil {
		log.Printf("Error listing parts: %s", err)

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	for _, p := range parts {
		res = append(res, mapper.ToProto(p))
	}

	return &inventoryv1.ListPartsResponse{
		Parts: res,
	}, nil
}
