package handler

import (
	"math/rand"
	"time"

	"ginfra/datasource"
	mw "ginfra/middleware"
	"ginfra/models"

	"github.com/gin-gonic/gin"
)

func DBTimedHandler(c *gin.Context) {

	// get the underlying request context
	ctx := c.Request.Context()

	db, _ := datasource.GormWithContext(ctx)
	var post models.Post

	mw.Logger(c).Info("begin sql...")
	rand.Seed(time.Now().UnixNano())
	err := db.First(&post, "id = ?", rand.Intn(10000)).Error
	mw.Logger(c).Info("end sql...")
	c.String(200, err.Error())
	return
}
