package repository

import (
	"context"
	"devices-api-go/config"
	"devices-api-go/internal/domain"
)

// DeviceRepository defines database operations for the Device domain.
type DeviceRepository struct{}

// NewDeviceRepository creates a new instance of DeviceRepository.
func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

// Save inserts a new device resource into the database schema.
func (r *DeviceRepository) Save(ctx context.Context, device *domain.Device) error {
	query := `INSERT INTO devices (name, brand, state, creation_time) 
	          VALUES ($1, $2, $3, NOW()) RETURNING id, creation_time`

	err := config.DB.QueryRow(ctx, query, device.Name, device.Brand, device.State).
		Scan(&device.ID, &device.CreationTime)
	return err
}

// FindAll fetches a paginated subset of devices starting AFTER the provided cursor ID.
func (r *DeviceRepository) FindAll(ctx context.Context, limit int, cursor int64) ([]domain.Device, error) {
	query := `SELECT id, name, brand, state, creation_time FROM devices 
	          WHERE id > $1 ORDER BY id ASC LIMIT $2`

	rows, err := config.DB.Query(ctx, query, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreationTime); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// FindByID retrieves a single device record by its unique identifier key.
func (r *DeviceRepository) FindByID(ctx context.Context, id int64) (*domain.Device, error) {
	query := `SELECT id, name, brand, state, creation_time FROM devices WHERE id = $1`

	var d domain.Device
	err := config.DB.QueryRow(ctx, query, id).
		Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreationTime)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// FindByBrand fetches devices belonging to a specific brand starting after the provided cursor ID.
func (r *DeviceRepository) FindByBrand(ctx context.Context, brand string, limit int, cursor int64) ([]domain.Device, error) {
	query := `SELECT id, name, brand, state, creation_time FROM devices 
	          WHERE id > $1 AND LOWER(brand) = LOWER($2) ORDER BY id ASC LIMIT $3`

	rows, err := config.DB.Query(ctx, query, cursor, brand, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreationTime); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// FindByState fetches devices currently set to a specific state starting after the provided cursor ID.
func (r *DeviceRepository) FindByState(ctx context.Context, state domain.DeviceState, limit int, cursor int64) ([]domain.Device, error) {
	query := `SELECT id, name, brand, state, creation_time FROM devices 
	          WHERE id > $1 AND state = $2 ORDER BY id ASC LIMIT $3`

	rows, err := config.DB.Query(ctx, query, cursor, state, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreationTime); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// Update executes a full replacement query for an existing device record.
func (r *DeviceRepository) Update(ctx context.Context, device *domain.Device) error {
	query := `UPDATE devices SET name = $1, brand = $2, state = $3 WHERE id = $4`
	_, err := config.DB.Exec(ctx, query, device.Name, device.Brand, device.State, device.ID)
	return err
}

// Delete permanently removes a specific record matching the targeting ID key.
func (r *DeviceRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM devices WHERE id = $1`
	_, err := config.DB.Exec(ctx, query, id)
	return err
}
