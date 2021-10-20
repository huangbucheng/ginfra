package seewo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"ginfra/errcode"
	"ginfra/utils"
)

//SeeWoAccessTokenBody seewo开放平台AccessToken
type SeeWoAccessTokenBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"open_id"`
	Scope        string `json:"scope"`
}

//SeeWoAccessTokenResponse seewo开放平台AccessToken
type SeeWoAccessTokenResponse struct {
	Body    SeeWoAccessTokenBody `json:"body"`
	Code    string               `json:"code"`
	SubCode string               `json:"sub_code"`
	Msg     string               `json:"message"`
	SubMsg  string               `json:"sub_message"`
}

//GetSeeWoAccessToken 获取seewo开放平台AccessToken
func GetSeeWoAccessToken(appid, secret string, code string) (*SeeWoAccessTokenResponse, error) {
	url := "https://openapi.seewo.com/api/oauth2/access_token?" +
		fmt.Sprintf("app_id=%s&app_secret=%s&code=%s&grant_type=authorization_code",
			appid, secret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			fmt.Sprintf("获取希沃AccessToken异常:%s", err.Error()))
	}
	defer resp.Body.Close()

	var response SeeWoAccessTokenResponse
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			fmt.Sprintf("解析希沃AccessToken异常:%s", err.Error()))
	}

	return &response, nil
}

type SeeWoUserInfo struct {
	Gender    int    `json:"gender"`
	NickName  string `json:"nickName"`
	PhotoUrl  string `json:"photoURL"`
	Uid       string `json:"uid"`
	AccountId string `json:"accountId"`
	Email     string `json:"email"`
}

//GetSeeWoUserInfo seewo开放平台获取用户信息
func GetSeeWoUserInfo(appid, secret, access_token, openid string) (*SeeWoUserInfo, error) {
	requrl := "https://openapi.seewo.com/user-center/user-api/user-base"
	body := []byte(fmt.Sprintf("{\"uid\":\"%s\"}", openid))

	headers := map[string]string{
		"Content-Type":     "application/json",
		"x-sw-app-id":      appid,
		"x-sw-auth-token":  access_token,
		"x-sw-timestamp":   strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		"x-sw-sign":        "",
		"x-sw-sign-type":   "hmac",
		"x-sw-req-path":    "/user-center/user-api/user-base",
		"x-sw-content-md5": strings.ToUpper(hex.EncodeToString(utils.MD5(body))),
		"x-sw-version":     "2",
	}

	// sign
	toSignList := []string{"x-sw-app-id", "x-sw-auth-token", "x-sw-timestamp",
		"x-sw-sign-type", "x-sw-req-path", "x-sw-content-md5", "x-sw-version"}
	sort.Strings(toSignList)

	var rawstring string
	for _, key := range toSignList {
		value, ok := headers[key]
		if !ok || len(value) == 0 {
			continue
		}
		rawstring += key + value
	}
	//fmt.Println("RAW:", rawstring)

	b := utils.HmacMD5([]byte(rawstring), secret)
	headers["x-sw-sign"] = strings.ToUpper(hex.EncodeToString(b))
	headers["x-sw-req-path"] = url.QueryEscape(headers["x-sw-req-path"])

	// request
	body, err := utils.PostRequest(requrl, headers, body)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			fmt.Sprintf("获取希沃用户信息异常:%s", err.Error()))
	}

	//fmt.Println(string(body))

	var userInfo SeeWoUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, errcode.NewCustomError(errcode.ErrCodeInternalError,
			fmt.Sprintf("解析希沃用户信息异常:%s", err.Error()))
	}

	return &userInfo, nil
}
