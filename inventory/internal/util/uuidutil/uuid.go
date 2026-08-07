package uuidutil

import "github.com/google/uuid"

func GetNewUUID() (string, error) {
	uuidBytes, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uuidBytes.String(), nil
}
