package middleware

import (
	"ginfra/protocol"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//RequestId middleware
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestId := c.Request.Header.Get(protocol.CtxRequestID)

		// Create request id with UUID4
		if requestId == "" {
			u4 := uuid.New()
			requestId = u4.String()
		}

		// Expose it for use in the application
		c.Set(protocol.CtxRequestID, requestId)

		// Set X-Request-Id header
		c.Writer.Header().Set(protocol.CtxRequestID, requestId)
		c.Next()
	}
}
