package tencent

import (
	"ginfra/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

var (
	stsSecretId  string
	stsSecretKey string
)

type StsCredential struct {
	// token。token长度和绑定的策略有关，最长不超过4096字节。
	Token *string `json:"Token,omitempty" name:"Token"`

	// 临时证书密钥ID。最长不超过1024字节。
	TmpSecretId *string `json:"TmpSecretId,omitempty" name:"TmpSecretId"`

	// 临时证书密钥Key。最长不超过1024字节。
	TmpSecretKey *string `json:"TmpSecretKey,omitempty" name:"TmpSecretKey"`

	// 临时证书有效的时间，返回 Unix 时间戳，精确到秒
	ExpiredTime *uint64 `json:"ExpiredTime,omitempty" name:"ExpiredTime"`
}

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	stsSecretId = cfg.GetString("qcloud.SecretID")
	stsSecretKey = cfg.GetString("qcloud.SecretKey")
}

func GetFederationToken(name, policy string) (*StsCredential, error) {

	credential := common.NewCredential(stsSecretId, stsSecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sts.tencentcloudapi.com"
	client, _ := sts.NewClient(credential, "ap-guangzhou", cpf)

	request := sts.NewGetFederationTokenRequest()
	request.Name = common.StringPtr(name)
	request.Policy = common.StringPtr(policy)
	request.DurationSeconds = common.Uint64Ptr(1800)

	response, err := client.GetFederationToken(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var res StsCredential
	res.Token = response.Response.Credentials.Token
	res.TmpSecretId = response.Response.Credentials.TmpSecretId
	res.TmpSecretKey = response.Response.Credentials.TmpSecretKey
	res.ExpiredTime = response.Response.ExpiredTime
	return &res, nil
}
