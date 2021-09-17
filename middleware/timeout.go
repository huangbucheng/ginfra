package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

//Timeout middleware
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				// c.AbortWithStatus(http.StatusGatewayTimeout)
			}
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
