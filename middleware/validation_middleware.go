package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"

	"devices-api-go/model"
	"devices-api-go/repository"
)

type ValidationMiddleware struct {
	repo *repository.DeviceRepository
}

func NewValidationMiddleware(repo *repository.DeviceRepository) *ValidationMiddleware {
	return &ValidationMiddleware{repo: repo}
}

func (m *ValidationMiddleware) GuardDeviceRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if idParam == "" {
			c.Next()
			return
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			_ = c.Error(errors.New("Invalid ID format"))
			c.Abort()
			return
		}

		existing, err := m.repo.FindByID(c.Request.Context(), id)
		if err != nil {
			_ = c.Error(errors.New("Device not found"))
			c.Abort()
			return
		}

		method := c.Request.Method

		// Rule 1: In use devices cannot be deleted.
		if method == http.MethodDelete {
			if existing.State == model.InUse {
				_ = c.Error(errors.New(model.ErrDeviceInUseDelete))
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if method == http.MethodPut || method == http.MethodPatch {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var payload map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &payload); err != nil {
				_ = c.Error(errors.New("Malformed JSON payload"))
				c.Abort()
				return
			}

			// Rule 2: Creation time cannot be updated.
			if _, exists := payload["creationTime"]; exists {
				_ = c.Error(errors.New(model.ErrCreationTimeImmutable))
				c.Abort()
				return
			}

			// Rule 3: Name and brand properties cannot be updated if the device is in use.
			if existing.State == model.InUse {
				if method == http.MethodPut {
					newName, _ := payload["name"].(string)
					newBrand, _ := payload["brand"].(string)

					if !strings.EqualFold(existing.Name, newName) || !strings.EqualFold(existing.Brand, newBrand) {
						_ = c.Error(errors.New(model.ErrDeviceInUseLocked))
						c.Abort()
						return
					}
				} else if method == http.MethodPatch {
					_, nameExists := payload["name"]
					_, brandExists := payload["brand"]

					if nameExists || brandExists {
						_ = c.Error(errors.New(model.ErrDeviceInUseLocked))
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}
