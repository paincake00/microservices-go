package memory

import (
	"sync"

	"github.com/paincake00/microservices-go/order/internal/entity"
)

type InMemoryOrderStorage struct {
	mx     sync.RWMutex
	orders map[string]entity.Order
}

func NewInMemoryOrderStorage() *InMemoryOrderStorage {
	return &InMemoryOrderStorage{
		orders: make(map[string]entity.Order),
	}
}
