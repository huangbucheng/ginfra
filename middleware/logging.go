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
func ContextLogger(logger *log.LoggerWrap) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), log.CtxLoggerFields, make(map[string]string))

		ctxReqId, _ := c.Value(log.CtxRequestID).(string)
		log.SetFields(ctx, map[string]string{
			"ClientIP":  c.ClientIP(),
			"RequestID": ctxReqId,
		})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func Ginzap() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				log.Logger.WithGinContext(c).Error(e)
			}
		} else {
			log.Logger.WithGinContext(c).Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("etime", end.Format(time.RFC3339)),
				zap.Duration("latency", latency),
			)
		}
	}
}
