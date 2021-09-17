package tencent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// QQ 互联 ------
//QQConnectTokenResponse QQ互联平台登录Code验证结果
type QQConnectTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Error        int    `json:"error"`
	ErrorMsg     string `json:"error_description"`
}

//QQConnectToken QQ互联平台登录Code验证
func QQConnectToken(appid, secret, redirect_uri, jscode string) (*QQConnectTokenResponse, error) {
	host := "https://graph.qq.com/oauth2.0/token?" +
		"client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&" +
		"grant_type=authorization_code&fmt=json"
	resp, err := http.Get(fmt.Sprintf(host, appid, secret, redirect_uri, jscode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response QQConnectTokenResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != 0 {
		return nil, errors.New(response.ErrorMsg)
	}

	return &response, nil
}

//QQConnectOpenIDResponse QQ互联平台OpenID结果
type QQConnectOpenIDResponse struct {
	OpenID   string `json:"openid"`
	UnionID  string `json:"unionid"`
	Error    int    `json:"error"`
	ErrorMsg string `json:"error_description"`
}

//QQConnectOpenID QQ互联平台获取OpenID
func QQConnectOpenID(access_token string) (*QQConnectOpenIDResponse, error) {
	host := "https://graph.qq.com/oauth2.0/me?" +
		"access_token=%s&unionid=1&fmt=json"
	resp, err := http.Get(fmt.Sprintf(host, access_token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response QQConnectOpenIDResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != 0 {
		return nil, errors.New(response.ErrorMsg)
	}
	return &response, err
}

//QQConnectUserInfoResponse QQ互联平台用户信息
type QQConnectUserInfoResponse struct {
	NickName  string `json:"nickname"`
	Gender    string `json:"gender"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Year      string `json:"year"`
	FigureUrl string `json:"figureurl"`
	Error     int    `json:"ret"`
	ErrorMsg  string `json:"msg"`
}

//QQConnectUserInfo QQ互联平台获取用户信息
func QQConnectUserInfo(appid, access_token, openid string) (*QQConnectUserInfoResponse, error) {
	host := "https://graph.qq.com/user/get_user_info?" +
		"access_token=%s&oauth_consumer_key=%s&openid=%s"
	resp, err := http.Get(fmt.Sprintf(host, access_token, appid, openid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response QQConnectUserInfoResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != 0 {
		return nil, errors.New(response.ErrorMsg)
	}
	return &response, err
}
