package memory

import (
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
)

type PostgresOrderStorage struct {
	pgClient *pgclient.PostgresClient
}

func NewPostgresOrderStorage(pgClient *pgclient.PostgresClient) *PostgresOrderStorage {
	return &PostgresOrderStorage{
		pgClient: pgClient,
	}
}
