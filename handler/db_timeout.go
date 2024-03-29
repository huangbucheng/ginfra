package handler

import (
	"math/rand"
	"time"

	"ginfra/datasource"
	"ginfra/log"
	"ginfra/models"

	"github.com/gin-gonic/gin"
)

//DBTimedHandler 示例
func DBTimedHandler(c *gin.Context) {

	// get the underlying request context
	ctx := c.Request.Context()

	db, _ := datasource.GormWithContext(ctx)
	var post models.Post

	log.WithGinContext(c).Info("begin sql...")
	rand.Seed(time.Now().UnixNano())
	err := db.First(&post, "id = ?", rand.Intn(10000)).Error
	log.WithGinContext(c).Info("end sql...")
	c.String(200, err.Error())
	return
}
