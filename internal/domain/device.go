package domain

import (
	"context"
	"time"
)

// DeviceState represents the current operational status of a device resource.
type DeviceState string

const (
	// Available status indicates the device is ready to be assigned.
	Available DeviceState = "AVAILABLE"
	// InUse status indicates the device is currently active and locked.
	InUse DeviceState = "IN_USE"
	// Inactive status indicates the device is disabled or powered off.
	Inactive DeviceState = "INACTIVE"
)

// Device represents the core database entity schema for device records.
// The tags specify how fields map to JSON payloads and database columns.
type Device struct {
	ID           int64       `json:"id" db:"id"`
	Name         string      `json:"name" db:"name"`
	Brand        string      `json:"brand" db:"brand"`
	State        DeviceState `json:"state" db:"state"`
	CreationTime time.Time   `json:"creationTime" db:"creation_time"`
}

// DeviceDTO defines the Data Transfer Object schema for client request payloads.
// These tags are utilized by validators to enforce mandatory input requirements.
type DeviceDTO struct {
	Name  string      `json:"name" binding:"required"`
	Brand string      `json:"brand" binding:"required"`
	State DeviceState `json:"state" binding:"required"`
}

// CursorPageResult wraps the dynamic slice array content alongside pagination tokens.
type CursorPageResult struct {
	Content    []Device `json:"content"`
	NextCursor int64    `json:"nextCursor"`
	HasNext    bool     `json:"hasNext"`
}

// DeviceRepository defines the data store contract for device lifecycle events.
type DeviceRepository interface {
	Save(ctx context.Context, device *Device) error
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Device, error)
	FindAll(ctx context.Context, limit int, cursor int64) ([]Device, error)
	FindByBrand(ctx context.Context, brand string, limit int, cursor int64) ([]Device, error)
	FindByState(ctx context.Context, state DeviceState, limit int, cursor int64) ([]Device, error)
}

// DeviceService defines the core business operations contract for devices.
type DeviceService interface {
	Create(ctx context.Context, dto DeviceDTO) (*Device, error)
	FullUpdate(ctx context.Context, id int64, dto DeviceDTO) (*Device, error)
	PartialUpdate(ctx context.Context, id int64, dto map[string]interface{}) (*Device, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, limit int, cursor int64) ([]Device, error)
	GetByID(ctx context.Context, id int64) (*Device, error)
	GetByBrand(ctx context.Context, brand string, limit int, cursor int64) ([]Device, error)
	GetByState(ctx context.Context, state DeviceState, limit int, cursor int64) ([]Device, error)
}
