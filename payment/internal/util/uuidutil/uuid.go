package uuidutil

import "github.com/google/uuid"

type UuidGeneratorV7 struct{}

func NewUuidGeneratorV7() *UuidGeneratorV7 {
	return &UuidGeneratorV7{}
}

func (g *UuidGeneratorV7) NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}
