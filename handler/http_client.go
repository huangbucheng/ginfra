package handler

import (
	mw "ginfra/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func UseHttpClient(c *gin.Context) {
	resp, err := req.Get("http://localhost", req.Param{"name": "roc", "age": "22"})
	if err != nil {
		mw.Logger(c).Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	body, _ := resp.ToString()
	mw.Logger(c).Info(body)

	c.String(http.StatusOK, "success")
}
