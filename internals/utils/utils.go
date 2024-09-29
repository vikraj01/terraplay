package utils

import "github.com/google/uuid"

// GenerateUUID creates a new UUID for session IDs
func GenerateUUID() string {
	return uuid.New().String()
}
