package tencent

import (
	"fmt"

	"ginfra/config"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var (
	captchaSecretId  string
	captchaSecretKey string

	captchaAppId  uint64
	captchaAppKey string
)

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	captchaSecretId = cfg.GetString("qcloud.SecretID")
	captchaSecretKey = cfg.GetString("qcloud.SecretKey")

	captchaAppId = cfg.GetUint64("captcha.AppID")
	captchaAppKey = cfg.GetString("captcha.AppKey")
}

func DescribeCaptchaResult(ticket, randstr, clientIp string) error {

	credential := common.NewCredential(captchaSecretId, captchaSecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "captcha.tencentcloudapi.com"
	client, _ := captcha.NewClient(credential, "", cpf)

	request := captcha.NewDescribeCaptchaResultRequest()

	request.CaptchaType = common.Uint64Ptr(9)
	// 前端回调函数返回的用户验证票据
	request.Ticket = common.StringPtr(ticket)
	request.UserIp = common.StringPtr(clientIp)
	// 前端回调函数返回的随机字符串
	request.Randstr = common.StringPtr(randstr)
	request.CaptchaAppId = common.Uint64Ptr(captchaAppId)
	request.AppSecretKey = common.StringPtr(captchaAppKey)

	response, err := client.DescribeCaptchaResult(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	if *response.Response.CaptchaCode != 1 {
		return fmt.Errorf(*response.Response.CaptchaMsg)
	}
	return nil
}
