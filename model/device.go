package model

import (
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
