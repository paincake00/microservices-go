package inventory

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func (s *Service) ListParts(ctx context.Context, partUuids []string) ([]*entity.Part, error) {
	resp, err := s.grpcInventoryClient.ListParts(
		ctx, &inventoryv1.ListPartsRequest{
			Filter: &inventoryv1.PartsFilter{
				Uuids: partUuids,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	parts := resp.GetParts()

	res := make([]*entity.Part, 0, len(parts))
	for _, part := range parts {
		res = append(
			res, &entity.Part{
				Price: part.GetPrice(),
			},
		)
	}

	return res, nil
}
