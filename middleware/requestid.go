package middleware

import (
	"ginfra/log"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestId := c.Request.Header.Get(log.CtxRequestID)

		// Create request id with UUID4
		if requestId == "" {
			u4 := uuid.NewV4()
			requestId = u4.String()
		}

		// Expose it for use in the application
		c.Set(log.CtxRequestID, requestId)

		// Set X-Request-Id header
		c.Writer.Header().Set(log.CtxRequestID, requestId)
		c.Next()
	}
}

func GetRequestId(c *gin.Context) string {
	if ctxReqId, ok := c.Value(log.CtxRequestID).(string); ok {
		return ctxReqId
	}

	return ""
}
