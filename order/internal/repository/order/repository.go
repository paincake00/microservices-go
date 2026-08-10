package order

import (
	"sync"

	"github.com/paincake00/microservices-go/order/internal/entity"
)

type InMemoryStorage struct {
	mx     sync.RWMutex
	orders map[string]entity.Order
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		orders: make(map[string]entity.Order),
	}
}
