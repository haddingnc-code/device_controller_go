package controller

import (
	"context"
	"devices-api-go/internal/domain"
	"devices-api-go/internal/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockHTTPService dynamically implements the domain.DeviceService interface for HTTP testing.
type MockHTTPService struct {
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockHTTPService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// Fulfill the rest of the contract with methods matching your specific service signatures.
func (m *MockHTTPService) Create(ctx context.Context, dto domain.DeviceDTO) (*domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) FullUpdate(ctx context.Context, id int64, dto domain.DeviceDTO) (*domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) PartialUpdate(ctx context.Context, id int64, dto map[string]interface{}) (*domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) GetAll(ctx context.Context, limit int, cursor int64) ([]domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) GetByID(ctx context.Context, id int64) (*domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) GetByBrand(ctx context.Context, brand string, limit int, cursor int64) ([]domain.Device, error) {
	return nil, nil
}
func (m *MockHTTPService) GetByState(ctx context.Context, state domain.DeviceState, limit int, cursor int64) ([]domain.Device, error) {
	return nil, nil
}

// TestController_Delete_ShouldReturnBadRequestWhenInUse verifies HTTP error mapping.
func TestController_Delete_ShouldReturnBadRequestWhenInUse(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockSvc := &MockHTTPService{
		DeleteFunc: func(ctx context.Context, id int64) error {
			return errors.New(domain.ErrDeviceInUseDelete)
		},
	}

	ctrl := NewDeviceController(mockSvc)

	router := gin.New()

	// Activate your global error handler so Gin intercepts c.Error(err)
	router.Use(middleware.GlobalErrorHandler())

	router.DELETE("/api/v1/devices/:id", ctrl.Delete)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/devices/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
