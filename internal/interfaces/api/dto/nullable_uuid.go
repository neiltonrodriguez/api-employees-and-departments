package dto

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

// NullableUUID is a custom type that handles UUID that can be null or empty string
type NullableUUID struct {
	UUID  uuid.UUID
	Valid bool // Valid is true if UUID is not NULL or empty
}

// UnmarshalJSON implements json.Unmarshaler
// Accepts: null, "", or a valid UUID string
func (nu *NullableUUID) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		nu.Valid = false
		nu.UUID = uuid.Nil
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Handle empty string as null
	if str == "" {
		nu.Valid = false
		nu.UUID = uuid.Nil
		return nil
	}

	// Parse UUID
	parsed, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	nu.UUID = parsed
	nu.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler
func (nu NullableUUID) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nu.UUID.String())
}

// Value implements driver.Valuer for database operations
func (nu NullableUUID) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}
	return nu.UUID, nil
}

// ToUUIDPointer converts NullableUUID to *uuid.UUID for compatibility
func (nu NullableUUID) ToUUIDPointer() *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}

// FromUUIDPointer creates a NullableUUID from *uuid.UUID
func FromUUIDPointer(u *uuid.UUID) NullableUUID {
	if u == nil {
		return NullableUUID{Valid: false, UUID: uuid.Nil}
	}
	return NullableUUID{Valid: true, UUID: *u}
}
