package handler

import (
	"fmt"
	"strconv"
	"time"

	"ginfra/config"
	"ginfra/errcode"
	"ginfra/log"
	"ginfra/protocol"
	"ginfra/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//GetDiscuzTokenRequest 获取discuz token的请求参数
type GetDiscuzTokenRequest struct {
}

//GetDiscuzTokenResponse 获取discuz token的响应参数
type GetDiscuzTokenResponse struct {
	Token string
}

//GetDiscuzToken 获取Discuz Token
func GetDiscuzToken(c *gin.Context) {
	var req GetDiscuzTokenRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithGinContext(c).Error(err.Error(), zap.String("error", errcode.ErrInvalidParam))
		protocol.SetErrResponse(c, protocol.ErrCodeInvalidParameter)
		return
	}

	claims, err := getClaimDataFromContext(c)
	if err != nil {
		protocol.SetErrResponse(c, err)
		return
	}

	cfg, _ := config.Parse("")
	privateKey := cfg.GetString("discuz.PrivateKey")
	if len(privateKey) == 0 {
		log.WithGinContext(c).Error(err.Error())
		protocol.SetErrResponse(c,
			errcode.NewCustomError(errcode.ErrCodeInternalError,
				"暂不支持签发Discuz Token"))
		return
	}

	// get discuz uid
	// discuz register: http://127.0.0.1:8090/apiv3/users/username.register
	//	POST Content-Type: application/json; charset=utf-8
	// {"username":"abc","password":"xxxx","nickname":"bob","passwordConfirmation":"xxxx","captchaRandStr":"","captchaTicket":"","code":""}

	data := make(map[string]interface{})
	data["sub"] = strconv.FormatUint(claims.Uid, 10)
	data["jti"] = "4a4550d0d9b3587c4f472038780452a3b17fd863c5aab7d14cca93037d49332726ab80dcbd9ddd59"
	data["aud"] = ""
	data["scopes"] = []interface{}{nil}
	data["exp"] = time.Now().Add(28 * 24 * time.Hour).Unix()
	data["iat"] = time.Now().Unix()
	data["nbf"] = time.Now().Unix()

	token, err := utils.CreateJWTTokenFromMapWithRS256(
		[]byte(privateKey),
		data,
	)
	if err != nil {
		protocol.SetErrResponse(c,
			errcode.NewCustomError(errcode.ErrCodeInternalError,
				fmt.Sprintf("create token error:%s", err.Error())))
		return
	}

	var resp GetDiscuzTokenResponse
	resp.Token = fmt.Sprintf("Bearer %s", token)

	protocol.SetResponseData(c, &resp)
}
