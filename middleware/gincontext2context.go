package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type gincontext_key string

var _gincontext_key gincontext_key = "GinContextKey"

//GinContextToContextMiddleware 中间件-将gin context存储到context中
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), _gincontext_key, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

//GinContextFromContext 中间件-从context中获取gin context
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(_gincontext_key)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
