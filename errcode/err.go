package errcode

import "fmt"

//CustomError 自定义Error类型
type CustomError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	error
}

//Error 自定义类型Error实现error interface
func (e CustomError) Error() string {
	return e.Message
}

//Wrap wrap inner err
func (e *CustomError) Wrap(err error) *CustomError {
	return &CustomError{
		Code:    e.Code,
		Message: fmt.Sprintf("%s->%s", e.Message, err.Error()),
	}
}

//Set set err message
func (e *CustomError) Set(err error) *CustomError {
	return &CustomError{
		Code:    e.Code,
		Message: err.Error(),
	}
}

// error code that can be assocation custom messages
var ErrCodeInternalError = "InternalError"
var ErrInvalidParam = "InvalidParameter"
var ErrNoAuthToken = "NoAuthToken"
var ErrInvalidAuthToken = "InvalidAuthToken"
var ErrExpiredAuthToken = "ExpiredAuthToken"
var ErrNoJWTClaims = "NoJWTClaims"

//NewCustomError 新建自定义Error
func NewCustomError(code, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

//ErrorCode 获取Error的ErrorCode
func ErrorCode(err error) string {
	pcerr, ok := err.(*CustomError)
	if ok {
		return pcerr.Code
	}
	cerr, ok := err.(CustomError)
	if ok {
		return cerr.Code
	}
	return ErrCodeInternalError
}
