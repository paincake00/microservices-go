package memory

import (
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
)

type PostgresOrderStorage struct {
	pgClient *pgclient.Client
}

func NewPostgresOrderStorage(pgClient *pgclient.Client) *PostgresOrderStorage {
	return &PostgresOrderStorage{
		pgClient: pgClient,
	}
}
