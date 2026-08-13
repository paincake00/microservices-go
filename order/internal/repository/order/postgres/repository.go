package memory

import (
	"sync"

	"github.com/paincake00/microservices-go/order/pkg/pgclient"
)

type PostgresOrderStorage struct {
	mx       sync.RWMutex
	pgClient *pgclient.PostgresClient
}

func NewPostgresOrderStorage(pgClient *pgclient.PostgresClient) *PostgresOrderStorage {
	return &PostgresOrderStorage{
		pgClient: pgClient,
	}
}
