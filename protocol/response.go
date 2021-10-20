package protocol

import (
	"encoding/json"
	"net/http"
	"time"

	"ginfra/errcode"
	"ginfra/utils"

	"github.com/gin-gonic/gin"
)

// 正常Response
type Response struct {
	Response innerResponse `json:"Response"`
}
type innerResponse struct {
	RequestId string `json:"RequestId"`
	Timestamp int64
	Data      interface{} `json:"Data,omitempty"`
}

// 错误返回 response，可选择使用
type ErrorResponse struct {
	Response innerErrorResponse `json:"Response"`
}

// 在response中封装错误的结构体，可选择使用
type innerErrorResponse struct {
	Error     errcode.CustomError `json:"Error"`
	RequestId string              `json:"RequestId"`
	Timestamp int64
}

//SetResponse 设置gin的response, for response without data field
func SetResponse(c *gin.Context, data interface{}) {
	var innerResp map[string]interface{}
	r := &innerResponse{
		RequestId: GetRequestId(c),
		Timestamp: time.Now().Unix(),
	}

	ja, _ := json.Marshal(r)
	json.Unmarshal(ja, &innerResp)

	jb, _ := json.Marshal(data)
	utils.DecodeJsonWithInt64(jb, &innerResp)

	resp := make(map[string]interface{})
	resp["Response"] = innerResp

	c.Set(CtxResponseCode, "OK")
	c.JSON(http.StatusOK, resp)
}

//SetResponseData 设置gin的response
func SetResponseData(c *gin.Context, data interface{}) {
	r := &innerResponse{
		RequestId: GetRequestId(c),
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	c.Set(CtxResponseCode, "OK")
	c.JSON(http.StatusOK, r)
}

//SetErrResponse 设置gin的error response
func SetErrResponse(c *gin.Context, err error) {
	cserr, ok := err.(*errcode.CustomError)
	if !ok {
		e, ok := err.(errcode.CustomError)
		if !ok {
			cserr = errcode.NewCustomError(errcode.ErrCodeInternalError, err.Error())
		} else {
			cserr = &e
		}
	}
	r := &ErrorResponse{
		Response: innerErrorResponse{
			RequestId: GetRequestId(c),
			Timestamp: time.Now().Unix(),
			Error:     *cserr,
		},
	}
	c.Set(CtxResponseCode, cserr.Code)
	c.JSON(http.StatusOK, r)
}
