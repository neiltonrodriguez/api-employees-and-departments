package errors

import "errors"

// Domain errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrValidation         = errors.New("validation error")
	ErrConflict           = errors.New("conflict error")
	ErrInvalidInput       = errors.New("invalid input")
	ErrManagerNotInDept   = errors.New("manager must be linked to the same department")
	ErrCycleDetected      = errors.New("cycle detected in department hierarchy")
	ErrInvalidCPF         = errors.New("invalid CPF")
	ErrDuplicateCPF       = errors.New("CPF already exists")
	ErrDuplicateRG        = errors.New("RG already exists")
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ConflictError represents a conflict error (uniqueness violation)
type ConflictError struct {
	Message string
	Field   string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
