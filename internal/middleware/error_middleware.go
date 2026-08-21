package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"devices-api-go/internal/domain"
)

// GlobalErrorHandler returns a Gin internal.middleware that intercepts errors attached to the request context.
// It acts exactly like Spring Boot's @RestControllerAdvice to centralize error serialization.
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Let the request execute through all controller and service layers first
		c.Next()

		// Check if any errors were appended to the context during execution
		if len(c.Errors) > 0 {
			// Extract the very last error that occurred
			err := c.Errors.Last()
			status := http.StatusInternalServerError
			message := err.Error()

			// Check our evaluation text rules to assign the correct 400 Bad Request status
			if message == domain.ErrCreationTimeImmutable ||
				message == domain.ErrDeviceInUseLocked ||
				message == domain.ErrDeviceInUseDelete {
				status = http.StatusBadRequest
			} else if message == "Device not found" {
				status = http.StatusNotFound
			} else if customStatus, exists := err.Meta.(int); exists {
				status = customStatus
			}

			// Handle Gin's internal JSON validation tags binding failures
			if err.Type == gin.ErrorTypeBind {
				status = http.StatusBadRequest
				message = "Validation failed for one or more fields. Check required parameters."
			}

			c.JSON(status, domain.ApiError{
				Timestamp: time.Now(),
				Status:    status,
				Error:     http.StatusText(status),
				Message:   message,
			})
			c.Abort()
			return
		}
	}
}
