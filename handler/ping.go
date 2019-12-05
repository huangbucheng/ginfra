package handler

import (
	mw "ginfra/middleware"
	"ginfra/utils"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	// set custom fileds into logger
	mw.SetLoggerField(c, map[string]string{
		mw.CtxProductID: "cbd271dec6133d7065bb5391a105f6ea",
		mw.CtxUserID:    "0qkkoqm22idmnmsno203u4nljdsf9",
	})

	utils.SearchCountInc(8000100)

	// mw.Log(c) returns the zap logger with custom fields
	mw.Logger(c).Info("ping...pong")
	c.String(200, "pong")
}
