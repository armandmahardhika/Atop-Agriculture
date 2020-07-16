package utils

import "github.com/google/uuid"

// GetGetUniqueID get uuid
func GetUniqueID() string {
	id := uuid.New()
	return id.String()
}
