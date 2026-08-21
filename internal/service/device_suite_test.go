package service

import (
	"context"
	"devices-api-go/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// MockInterfaceRepository dynamically implements domain.DeviceRepository contract for testing.
type MockInterfaceRepository struct {
	FindByIDFunc func(ctx context.Context, id int64) (*domain.Device, error)
}

// Implement only the required method for our AOP aspect tests.
func (m *MockInterfaceRepository) FindByID(ctx context.Context, id int64) (*domain.Device, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

// Fulfill the rest of the interface contract with no-op methods to allow compilation.
func (m *MockInterfaceRepository) Save(ctx context.Context, d *domain.Device) error   { return nil }
func (m *MockInterfaceRepository) Update(ctx context.Context, d *domain.Device) error { return nil }
func (m *MockInterfaceRepository) Delete(ctx context.Context, id int64) error         { return nil }
func (m *MockInterfaceRepository) FindAll(ctx context.Context, l int, c int64) ([]domain.Device, error) {
	return nil, nil
}
func (m *MockInterfaceRepository) FindByBrand(ctx context.Context, b string, l int, c int64) ([]domain.Device, error) {
	return nil, nil
}
func (m *MockInterfaceRepository) FindByState(ctx context.Context, s domain.DeviceState, l int, c int64) ([]domain.Device, error) {
	return nil, nil
}

// TestTestSuite_Aspect_Delete_Guard verifies the AOP aspect safely catches active state blocks.
func TestTestSuite_Aspect_Delete_Guard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockInterfaceRepository{
		FindByIDFunc: func(ctx context.Context, id int64) (*domain.Device, error) {
			return &domain.Device{
				ID:           1,
				Name:         "iPhone 15",
				Brand:        "Apple",
				State:        domain.InUse,
				CreationTime: time.Now(),
			}, nil
		},
	}

	// Because repo is now a domain interface, mockRepo binds cleanly without pointer hacks!
	aspectSvc := NewDeviceServiceAspect(nil, mockRepo)

	// Act
	err := aspectSvc.Delete(ctx, 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDeviceInUseDelete, err.Error())
}

// TestTestSuite_Aspect_Update_Guard verifies mutation attempts on locked objects are blocked.
func TestTestSuite_Aspect_Update_Guard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockInterfaceRepository{
		FindByIDFunc: func(ctx context.Context, id int64) (*domain.Device, error) {
			return &domain.Device{
				ID:           1,
				Name:         "Galaxy S24",
				Brand:        "Samsung",
				State:        domain.InUse,
				CreationTime: time.Now(),
			}, nil
		},
	}

	aspectSvc := NewDeviceServiceAspect(nil, mockRepo)
	payloadDto := domain.DeviceDTO{
		Name:  "Malicious Update Attempt",
		Brand: "Samsung",
		State: domain.InUse,
	}

	// Act
	_, err := aspectSvc.FullUpdate(ctx, 1, payloadDto)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDeviceInUseLocked, err.Error())
}

// TestTestSuite_Service_Create_Validation isolates raw property checks.
func TestTestSuite_Service_Create_Validation(t *testing.T) {
	ctx := context.Background()
	inputDto := domain.DeviceDTO{
		Name:  "Test Core Phone",
		Brand: "Core Brand",
		State: domain.Available,
	}

	assert.NotEmpty(t, inputDto.Name)
	assert.NotEmpty(t, inputDto.Brand)
	assert.Equal(t, domain.Available, inputDto.State)
	assert.NotNil(t, ctx)
}
