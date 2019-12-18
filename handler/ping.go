package handler

import (
	"ginfra/log"
	"ginfra/utils"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	// set custom fileds into logger
	log.SetFieldsByGin(c, map[string]string{
		log.CtxProductID: "cbd271dec6133d7065bb5391a105f6ea",
		log.CtxUserID:    "0qkkoqm22idmnmsno203u4nljdsf9",
	})

	utils.SearchCountInc(8000100)

	// mw.Log(c) returns the zap logger with custom fields
	log.Logger.WithGinContext(c).Info("ping...pong")
	c.String(200, "pong")
}
