package protocol

import "github.com/gin-gonic/gin"

var CtxRequestID = "X-Request-Id"
var CtxClientIP = "X-Forwarded-For"

var CtxProductID = "X-Product-ID"       // 机器人ID
var CtxLoggerFields = "X-Logger-Fields" // 自定义日志字段集，map[string]string

// 自定义日志字段
var CtxUserID = "X-User-ID"             // 用户ID, 问答用户
var CtxCustomerID = "X-Customer-ID"     // 客户ID, 商户、客服等
var CtxResponseCode = "X-Response-Code" // 返回码

//GetUserId 获取gin请求UserId
func GetUserId(c *gin.Context) string {
	if ctxUserId, ok := c.Value(CtxUserID).(string); ok {
		return ctxUserId
	}

	return ""
}

//GetRequestId 获取gin请求ID
func GetRequestId(c *gin.Context) string {
	if ctxReqId, ok := c.Value(CtxRequestID).(string); ok {
		return ctxReqId
	}

	return ""
}

//GetResponseCode 获取gin响应Code
func GetResponseCode(c *gin.Context) string {
	if ctxRespCode, ok := c.Value(CtxResponseCode).(string); ok {
		return ctxRespCode
	}

	return ""
}

//GetClientIP 获取gin ClientIP
func GetClientIP(c *gin.Context) string {
	if ctxClientIP, ok := c.Value(CtxClientIP).(string); ok {
		return ctxClientIP
	}

	return ""
}
