package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		// Send it back to client
		c.Writer.Header().Set("X-Request-ID", id)
		// Make it available in context for logging, tracing, etc.
		c.Set("request_id", id)
		c.Next()
	}
}
