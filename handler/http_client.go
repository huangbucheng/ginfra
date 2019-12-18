package handler

import (
	"net/http"

	"ginfra/log"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func UseHttpClient(c *gin.Context) {
	resp, err := req.Get("http://localhost", req.Param{"name": "roc", "age": "22"})
	if err != nil {
		log.Logger.WithGinContext(c).Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	body, _ := resp.ToString()
	log.Logger.WithGinContext(c).Info(body)

	c.String(http.StatusOK, "success")
}
