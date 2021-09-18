package tencent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

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
}

//QQDocsToken 腾讯文档开放平台授权验证
func QQDocsToken(appid, secret, redirect_uri, code string) (*QQDocsTokenResponse, error) {
	host := "https://docs.qq.com/oauth/v2/token?" +
		"client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&" +
		"grant_type=authorization_code"
	body, err := utils.GetRequest(fmt.Sprintf(host, appid, secret, redirect_uri, code), nil)
	if err != nil {
		return nil, err
	}

	var response QQDocsTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

//QQDocsTempUrlResponse 腾讯文档临时URL
type QQDocsTempUrlResponse struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		TempUrl string `json:"tempURL"`
	} `json:"data"`
}

//QueryQQDocsTempUrl 获取腾讯文档临时URL
func QueryQQDocsTempUrl(appid, token, openid, fileid string) (*QQDocsTempUrlResponse, error) {
	host := "https://docs.qq.com/openapi/drive/v2/util/temp-url?" +
		"fileID=%s&type=open"
	headers := map[string]string{
		"Access-Token": token,
		"Client-Id":    appid,
		"Open-Id":      openid,
		"Content-Type": "application/x-www-form-urlencoded",
	}
	body, err := utils.PostRequest(fmt.Sprintf(host, url.QueryEscape(fileid)), headers, nil)
	if err != nil {
		return nil, err
	}

	var response QQDocsTempUrlResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 0 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}

//QQDocsConverterResponse 腾讯文档ID转换
type QQDocsConverterResponse struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		FileID    string `json:"fileID"`
		EncodedID string `json:"encodedID"`
	} `json:"data"`
}

//QQDocsConverter 腾讯文档ID转换
func QQDocsConverter(appid, token, openid, srcID string, convType int) (*QQDocsConverterResponse, error) {
	host := "https://docs.qq.com/openapi/drive/v2/util/converter?" +
		"value=%s&type=%d"
	headers := map[string]string{
		"Access-Token": token,
		"Client-Id":    appid,
		"Open-Id":      openid,
		"Content-Type": "application/x-www-form-urlencoded",
	}
	body, err := utils.GetRequest(fmt.Sprintf(host, url.QueryEscape(srcID), convType), headers)
	if err != nil {
		return nil, err
	}

	var response QQDocsConverterResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 0 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}
