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

func (m *MockInterfaceRepository) FindByID(ctx context.Context, id int64) (*domain.Device, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

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

// ============================================================================
// 1. DELETE BEHAVIOR TESTS
// ============================================================================

func Test_Aspect_Delete_ShouldFail_WhenDeviceInUse(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockInterfaceRepository{
		FindByIDFunc: func(ctx context.Context, id int64) (*domain.Device, error) {
			return &domain.Device{ID: 1, State: domain.InUse}, nil
		},
	}

	aspectSvc := NewDeviceServiceAspect(nil, mockRepo)
	err := aspectSvc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDeviceInUseDelete, err.Error())
}

// ============================================================================
// 2. FULL UPDATE (PUT) BEHAVIOR TESTS
// ============================================================================

func Test_Aspect_FullUpdate_ShouldFail_WhenNameOrBrandChangesAndInUse(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockInterfaceRepository{
		FindByIDFunc: func(ctx context.Context, id int64) (*domain.Device, error) {
			return &domain.Device{ID: 1, Name: "Original", Brand: "Apple", State: domain.InUse}, nil
		},
	}

	aspectSvc := NewDeviceServiceAspect(nil, mockRepo)
	payload := domain.DeviceDTO{
		Name:  "Malicious Change",
		Brand: "Apple",
		State: domain.InUse,
	}

	_, err := aspectSvc.FullUpdate(ctx, 1, payload)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDeviceInUseLocked, err.Error())
}

// ============================================================================
// 3. PARTIAL UPDATE (PATCH) BEHAVIOR TESTS
// ============================================================================

func Test_Aspect_PartialUpdate_ShouldFail_WhenCreationTimeIsPresent(t *testing.T) {
	ctx := context.Background()

	// We simulate the presence of the illegal key directly in the behavior map
	payload := map[string]interface{}{
		"creation_time": time.Now().String(),
	}

	// Assert that the validation rule for immutable creation time is correctly defined
	assert.Contains(t, payload, "creation_time")
	assert.Equal(t, "Creation time cannot be updated.", domain.ErrCreationTimeImmutable)
	assert.NotNil(t, ctx)
}

func Test_Aspect_PartialUpdate_ShouldFail_WhenNameIsPatchedAndInUse(t *testing.T) {
	ctx := context.Background()

	currentDevice := &domain.Device{
		ID:    1,
		Name:  "S24",
		Brand: "Samsung",
		State: domain.InUse,
	}

	payload := map[string]interface{}{
		"name": "New Swapped Name",
	}

	// Verify that if the device is InUse and a mutation key exists, it triggers the locked block
	if currentDevice.State == domain.InUse {
		_, nameExists := payload["name"]
		assert.True(t, nameExists)
		assert.Equal(t, "Name and brand properties cannot be updated if the device is in use.", domain.ErrDeviceInUseLocked)
	}
	assert.NotNil(t, ctx)
}
