package service

import (
	"context"
	"devices-api-go/model"
	"devices-api-go/repository"
)

// DeviceService orchestrates data flows between controllers and repositories.
type DeviceService struct {
	repo *repository.DeviceRepository
}

// NewDeviceService acts as a constructor initializing the service with its repository dependency.
func NewDeviceService(repo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		repo: repo,
	}
}

// Create handles the setup of a new device entity and saves it.
func (s *DeviceService) Create(ctx context.Context, dto model.DeviceDTO) (*model.Device, error) {
	device := &model.Device{
		Name:  dto.Name,
		Brand: dto.Brand,
		State: dto.State,
	}

	err := s.repo.Save(ctx, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// GetAll fetches a paginated collection block of all existing devices.
func (s *DeviceService) GetAll(ctx context.Context, limit, offset int) ([]model.Device, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

// GetByID locates a single device record or propagates an error if missing.
func (s *DeviceService) GetByID(ctx context.Context, id int64) (*model.Device, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByBrand handles filtered pagination based on a specific brand string.
func (s *DeviceService) GetByBrand(ctx context.Context, brand string, limit, offset int) ([]model.Device, error) {
	return s.repo.FindByBrand(ctx, brand, limit, offset)
}

// GetByState handles filtered pagination based on an operational state enum.
func (s *DeviceService) GetByState(ctx context.Context, state model.DeviceState, limit, offset int) ([]model.Device, error) {
	return s.repo.FindByState(ctx, state, limit, offset)
}

// FullUpdate overwrites an existing record with clean validated payload properties.
func (s *DeviceService) FullUpdate(ctx context.Context, id int64, dto model.DeviceDTO) (*model.Device, error) {
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	device.Name = dto.Name
	device.Brand = dto.Brand
	device.State = dto.State

	err = s.repo.Update(ctx, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// PartialUpdate modifies only the specific payload key-value fields sent by the client.
func (s *DeviceService) PartialUpdate(ctx context.Context, id int64, updates map[string]interface{}) (*model.Device, error) {
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if val, exists := updates["name"]; exists {
		device.Name = val.(string)
	}
	if val, exists := updates["brand"]; exists {
		device.Brand = val.(string)
	}
	if val, exists := updates["state"]; exists {
		device.State = model.DeviceState(val.(string))
	}

	err = s.repo.Update(ctx, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// Delete permanently removes a targeted device record from the data tier.
func (s *DeviceService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
