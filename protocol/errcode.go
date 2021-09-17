package protocol

import (
	"ginfra/errcode"
)

var ErrCodeInvalidParameter *errcode.CustomError = &errcode.CustomError{
	Code:    errcode.ErrInvalidParam,
	Message: "请求参数解析错误",
}

var ErrCodeMissingParameter *errcode.CustomError = &errcode.CustomError{
	Code:    "MissingParam",
	Message: "缺失必填请求参数",
}

var ErrCodeInvalidClaims *errcode.CustomError = &errcode.CustomError{
	Code:    "InvalidJWTClaims",
	Message: "用户登录态信息无效，请重新登录",
}

var ErrCodeUnAuthorized *errcode.CustomError = &errcode.CustomError{
	Code:    "UnauthorizedOperation",
	Message: "无权限，请检查账号是否有权限访问相关数据",
}

var ErrCodeInvalidWXCode *errcode.CustomError = &errcode.CustomError{
	Code:    "InvalidWXCode",
	Message: "微信登录CODE无效",
}

var ErrCodeDBException *errcode.CustomError = &errcode.CustomError{
	Code:    "DataException",
	Message: "操作数据异常，请稍后重试",
}
