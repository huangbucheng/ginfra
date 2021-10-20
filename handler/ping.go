package handler

import (
	"ginfra/log"
	"ginfra/prom"
	"ginfra/protocol"

	"github.com/gin-gonic/gin"
)

//Ping 示例
func Ping(c *gin.Context) {
	// set custom fileds into logger
	log.Logger(c).Set(
		protocol.CtxProductID, "cbd271dec6133d7065bb5391a105f6ea").Set(
		protocol.CtxUserID, "0qkkoqm22idmnmsno203u4nljdsf9")

	prom.SearchCountInc(8000100)

	// mw.Log(c) returns the zap logger with custom fields
	log.WithGinContext(c).Info("ping...pong")
	c.String(200, "pong")
}
