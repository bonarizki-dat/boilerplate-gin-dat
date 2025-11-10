package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request.
//
// The request ID is:
// - Generated using UUID v4
// - Set in gin context as "request_id"
// - Added to response header as "X-Request-ID"
// - Can be used for request tracing and logging
//
// Usage in handlers:
//
//	requestID := c.GetString("request_id")
//	logger.WithField("request_id", requestID).Info("processing request")
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header (from load balancer/proxy)
		requestID := c.GetHeader("X-Request-ID")

		// Generate new UUID if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for handlers to use
		c.Set("request_id", requestID)

		// Add to response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
