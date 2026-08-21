package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"devices-api-go/middleware"
	"devices-api-go/model"
)

// setupTestRouter builds a minimal decoupled Gin engine to test middleware routing isolation.
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Register the global error handler configuration
	r.Use(middleware.GlobalErrorHandler())
	return r
}

func TestMiddleware_EnforceDeleteRestrictionWhenInUse(t *testing.T) {
	// 1. Initialize the mock test routing server environment
	r := setupTestRouter()

	// Create a dynamic route handler simulating the DELETE guard mapping behavior
	r.DELETE("/api/v1/devices/:id", func(c *gin.Context) {
		// Securely attach the error type matching your business constraints
		_ = c.Error(errors.New(model.ErrDeviceInUseDelete))
		c.Abort()
	})

	// 2. Perform the mock HTTP call execution loop
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/devices/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 3. Assert target status codes and evaluation text blocks match criteria
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "In use devices cannot be deleted.")
}

func TestMiddleware_EnforceCreationTimeImmutability(t *testing.T) {
	r := setupTestRouter()

	r.PATCH("/api/v1/devices/:id", func(c *gin.Context) {
		_ = c.Error(errors.New(model.ErrCreationTimeImmutable))
		c.Abort()
	})

	payload := map[string]interface{}{
		"creationTime": "2026-08-21T00:00:00Z",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/devices/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Creation time cannot be updated.")
}

func TestMiddleware_EnforceUpdateRestrictionWhenInUse(t *testing.T) {
	r := setupTestRouter()

	r.PUT("/api/v1/devices/:id", func(c *gin.Context) {
		_ = c.Error(errors.New(model.ErrDeviceInUseLocked))
		c.Abort()
	})

	dto := model.DeviceDTO{
		Name:  "iPhone 14 Pro",
		Brand: "Apple",
		State: model.InUse,
	}
	body, _ := json.Marshal(dto)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/devices/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Name and brand properties cannot be updated if the device is in use.")
}
