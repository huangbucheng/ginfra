package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"ginfra/errcode"
)

// 微信开放平台 ------
//WxCode2SessionResponse 微信开放平台登录Code验证结果
type WxCode2SessionResponse struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

//WxCode2Session 微信开放平台登录Code验证
func WxCode2Session(appid, secret, jscode string) (*WxCode2SessionResponse, error) {
	host := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	resp, err := http.Get(fmt.Sprintf(host, appid, secret, jscode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response WxCode2SessionResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)

	return &response, err
}

//WxOpenAccessTokenResponse 微信开放平台AccessToken
type WxOpenAccessTokenResponse struct {
	Access_Token  string `json:"access_token"`
	Expires_In    int    `json:"expires_in"`
	Refresh_Token string `json:"refresh_token"`
	Openid        string `json:"openid"`
	Scope         string `json:"scope"`
	Unionid       string `json:"unionid"`
	Errcode       int    `json:"errcode"`
	Errmsg        string `json:"errmsg"`
}

//GetWxOpenAccessToken 获取微信开放平台AccessToken
func GetWxOpenAccessToken(appid, secret string, code string) (*WxOpenAccessTokenResponse, error) {
	wxloginurl := "https://api.weixin.qq.com/sns/oauth2/access_token?" +
		fmt.Sprintf("appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			appid, secret, code)

	resp, err := http.Get(wxloginurl)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			"获取微信AccessToken异常，请重试")
	}
	defer resp.Body.Close()

	var wxresponse WxOpenAccessTokenResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &wxresponse)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			"获取微信AccessToken异常，请重试")
	}

	return &wxresponse, nil
}

//WxOpenUserInfo 微信开放平台获取的用户信息
type WxOpenUserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
	Errcode    int      `json:"errcode"`
	Errmsg     string   `json:"errmsg"`
}

//GetWxOpenUserInfo 微信开放平台获取用户信息
func GetWxOpenUserInfo(access_token, openid string) (*WxOpenUserInfo, error) {
	wxurl := "https://api.weixin.qq.com/sns/userinfo?" +
		fmt.Sprintf("access_token=%s&openid=%s", access_token, openid)

	resp, err := http.Get(wxurl)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			"获取微信用户信息异常，请重试")
	}
	defer resp.Body.Close()

	var wxresponse WxOpenUserInfo
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &wxresponse)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			"获取微信用户信息异常，请重试")
	}

	return &wxresponse, nil
}

// 微信公众号 ------
//WxOffiAcctAccessTokenResponse 微信公众号AccessToken
type WxOffiAcctAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

//GetWxOffiAcctAccessToken 获取微信公众号AccessToken
func GetWxOffiAcctAccessToken(appid, secret string) (*WxOffiAcctAccessTokenResponse, error) {
	host := "https://api.weixin.qq.com/cgi-bin/token?appid=%s&secret=%s&grant_type=client_credential"
	resp, err := http.Get(fmt.Sprintf(host, appid, secret))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response WxOffiAcctAccessTokenResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)

	return &response, err
}

//WxOffiAcctJSAPITicketResponse 微信公众号JS API Ticket
type WxOffiAcctJSAPITicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

//GetWxOffiAcctJSAPITicket 获取微信公众号JS API Ticket
func GetWxOffiAcctJSAPITicket(access_token string) (*WxOffiAcctJSAPITicketResponse, error) {
	host := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	resp, err := http.Get(fmt.Sprintf(host, access_token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response WxOffiAcctJSAPITicketResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)

	return &response, err
}

//GetWxOffiAcctJSSDKSignature 获取微信公众号JS SDK Signature
func GetWxOffiAcctJSSDKSignature(ticket, nonce, url string, ts int64) string {
	rawstring := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonce, ts, url)
	h := sha1.New()
	h.Write([]byte(rawstring))
	return hex.EncodeToString(h.Sum(nil))
}

//CheckWxOffiAcctSignature 微信公众号签名检查
func CheckWxOffiAcctSignature(signature, timestamp, nonce, token string) bool {
	arr := []string{timestamp, nonce, token}
	// 字典序排序
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return Sha1(b.String()) == signature
}

// 进行Sha1编码
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
