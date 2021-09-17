package tencent

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"ginfra/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	facefusion "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/facefusion/v20181201"
)

var (
	faceSecretId  string
	faceSecretKey string
)

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	faceSecretId = cfg.GetString("qcloud.SecretID")
	faceSecretKey = cfg.GetString("qcloud.SecretKey")
}

//ToBase64 to base64 string
func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

//ReadImage read image from local file
func ReadImage(filename string) string {
	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}

	return ToBase64(bytes)
}

//FaceFusion 腾讯云人脸融合接口
func FaceFusion(projId, moduleId, image string) (string, error) {

	credential := common.NewCredential(
		faceSecretId,
		faceSecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "facefusion.tencentcloudapi.com"
	client, _ := facefusion.NewClient(credential, "", cpf)

	request := facefusion.NewFaceFusionRequest()

	request.ProjectId = common.StringPtr(projId)
	request.ModelId = common.StringPtr(moduleId)
	request.Image = common.StringPtr(image)
	request.RspImgType = common.StringPtr("base64")

	response, err := client.FaceFusion(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", err
	}
	if err != nil {
		return "", err
	}
	fmt.Printf("%s", response.ToJsonString())
	if response.Response.Image == nil {
		return "", nil
	}
	return *response.Response.Image, nil
}
