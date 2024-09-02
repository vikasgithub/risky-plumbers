package entity

import "github.com/google/uuid"

// GenerateID Generates a new UUID
func GenerateID() string {
	return uuid.New().String()
}
