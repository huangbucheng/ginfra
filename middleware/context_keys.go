package middleware

var CtxLoggerKey = "X-Logger"
var CtxRequestID = "X-Request-Id"

var CtxProductID = "X-Product-ID"       // 机器人ID
var CtxLoggerFields = "X-Logger-Fields" // 自定义日志字段集，map[string]string

// 自定义日志字段
var CtxUserID = "X-User-ID"         // 用户ID, 问答用户
var CtxCustomerID = "X-Customer-ID" // 客户ID, 商户、客服等
