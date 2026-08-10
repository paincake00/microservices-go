package entity

import "errors"

var (
	ErrPartsNotFound = errors.New("parts not found")
	ErrOrderConflict = errors.New("order already paid")
)
