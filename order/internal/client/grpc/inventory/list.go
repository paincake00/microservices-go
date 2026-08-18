package inventory

import (
	"context"
	"log"
	"time"

	"github.com/paincake00/microservices-go/order/internal/entity"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

func (s *Service) ListParts(ctx context.Context, partUuids []string) ([]*entity.Part, error) {
	start := time.Now()
	log.Printf("ListParts: before inventory call")

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

	log.Printf("ListParts: after inventory call, took %s", time.Since(start))

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
