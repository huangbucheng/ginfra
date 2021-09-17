package handler

import (
	"net/http"

	"ginfra/datasource"
	"ginfra/models"

	"github.com/gin-gonic/gin"
)

//PostCreate 示例
func PostCreate(c *gin.Context) {
	// tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished

	post := &models.Post{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}

	db, _ := datasource.Gormv2(c.Request.Context())
	err := post.Insert(db)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Response": gin.H{
				"Error": gin.H{
					"Code":    500,
					"Message": err.Error(),
				},
				"post": post,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Response": gin.H{
			"Message": "success",
		},
	})
}

//PostGet 示例
func PostGet(c *gin.Context) {
	id := c.Param("id")
	db, _ := datasource.Gormv2(c.Request.Context())

	post, err := models.GetPostById(db, id)
	if err != nil || !post.IsPublished {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	post.View++
	post.UpdateView(db)
	c.JSON(http.StatusOK, gin.H{
		"Response": gin.H{
			"Message": "success",
			"Post":    post,
		},
	})
}
