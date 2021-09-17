package tencent

import (
	"encoding/json"
	"errors"
	"fmt"

	"ginfra/utils"
)

//QQDocsTokenResponse 腾讯文档开放平台授权验证结果
type QQDocsTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"user_id"`
	Scope        string `json:"scope"`
	Error        int    `json:"error"`
	ErrorMsg     string `json:"error_description"`
}

//QQDocsToken 腾讯文档开放平台授权验证
func QQDocsToken(appid, secret, redirect_uri, code string) (*QQDocsTokenResponse, error) {
	host := "https://docs.qq.com/oauth/v2/token?" +
		"client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&" +
		"grant_type=authorization_code"
	body, err := utils.GetRequest(fmt.Sprintf(host, appid, secret, redirect_uri, code))
	if err != nil {
		return nil, err
	}

	var response QQDocsTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != 0 {
		return nil, errors.New(response.ErrorMsg)
	}

	return &response, nil
}
