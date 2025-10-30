package uuid

import (
	"github.com/google/uuid"
)

// NewV7 generates a new UUIDv7
// UUIDv7 is time-ordered and suitable for database primary keys
func NewV7() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if v7 generation fails
		return uuid.New()
	}
	return id
}
