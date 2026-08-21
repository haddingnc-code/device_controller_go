package domain

import "time"

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
