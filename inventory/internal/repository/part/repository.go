package part

import (
	"sync"

	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

type Repository struct {
	mx    sync.RWMutex
	parts map[string]*model.Part
}

func NewRepository() *Repository {
	return &Repository{
		parts: make(map[string]*model.Part),
	}
}
