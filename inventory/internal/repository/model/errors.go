package model

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrSerializeMetadata   = errors.New("error serializing metadata")
	ErrUnserializeMetadata = errors.New("error unserializing metadata")
)
