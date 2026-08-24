package domain

import (
	"errors"
	"time"
)

// ErrDeviceNotFound is returned by the repository when no device matches the given ID.
// Using a sentinel error (checked with errors.Is) instead of a string comparison
// avoids depending on the exact wording of the underlying driver error.
var ErrDeviceNotFound = errors.New("device not found")

// ApiError defines the standardized JSON structure for all API error responses.
type ApiError struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
}

// Business rule constraints enforced by the validation layer.
const (
	ErrCreationTimeImmutable = "Creation time cannot be updated."
	ErrDeviceInUseLocked     = "Name and brand properties cannot be updated if the device is in use."
	ErrDeviceInUseDelete     = "In use devices cannot be deleted."
)
