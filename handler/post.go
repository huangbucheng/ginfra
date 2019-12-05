package handler

import (
	"ginfra/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	err := post.Insert(c.Request.Context())
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

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(c.Request.Context(), id)
	if err != nil || !post.IsPublished {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	post.View++
	post.UpdateView(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"Response": gin.H{
			"Message": "success",
			"Post":    post,
		},
	})
}
