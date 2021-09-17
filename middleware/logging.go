package middleware

import (
	"context"
	"fmt"
	"time"

	"ginfra/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinCustomLogFormat defines custom log format
func GinCustomLogFormat(param gin.LogFormatterParams) string {

	// your custom format
	requestId := ""
	if v, found := param.Keys["X-Request-Id"]; found {
		requestId = v.(string)
	}
	return fmt.Sprintf("%s | %d | %s | %s | %s | %s | %s | %s | %s\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		requestId,
		param.Request.Method,
		param.Path,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// ContextLogger middleware creates a derived logger to include logging of the
// Request ID, and inserts it into the context object
func ContextLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.NewContextLogger(logger)
		ctx := context.WithValue(c.Request.Context(), log.CtxLoggerKey, l)

		var client_ip string = c.ClientIP()
		if forwarded_for, ok := c.Request.Header["X-Forwarded-For"]; ok {
			if len(forwarded_for) > 0 {
				client_ip = forwarded_for[0]
			}
		}
		c.Set(log.CtxClientIP, client_ip)

		ctxReqId, _ := c.Value(log.CtxRequestID).(string)
		l.Set("ClientIP", client_ip).Set("RequestID", ctxReqId)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
