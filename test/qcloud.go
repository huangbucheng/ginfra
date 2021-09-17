package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"ginfra/utils"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func sendSms() {
	err := utils.SendSms("+8618500000000", "test", "100000", []string{"123456", "5"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sms Sent Successfully!")
}

func captcha() {
	err := utils.DescribeCaptchaResult("xxxx",
		"xxxx", "127.0.0.1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Verify Successfully!")
}

func sts() {
	policy := "{\"statement\":[{\"action\":[\"name/cos:PutObject\",\"name/cos:PostObject\",\"name/cos:InitiateMultipartUpload\",\"name/cos:UploadPart\",\"name/cos:CompleteMultipartUpload\",\"name/cos:AbortMultipartUpload\"],\"effect\":\"allow\",\"resource\":[\"qcs::cos:ap-guangzhou:uid/xxxx:test-bucket/*\"]}],\"version\":\"2.0\"}"
	resp, err := utils.GetFederationToken("dummy", policy)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func testcoswithTmpSecret() {
	// 使用STS临时密钥
	u, _ := url.Parse("https://test-bucket.cos.ap-guangzhou.myqcloud.com")
	fmt.Println(u)
	b := &cos.BaseURL{BucketURL: u}
	_CosClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 使用GetFederationToken接口返回的临时密钥
			SecretID:     "xxx",
			SecretKey:    "xxx",
			SessionToken: "xxx",
		},
	})

	var err error
	_, err = _CosClient.Object.PutFromFile(context.Background(), "image/test.jpg", "./test.jpg", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PutFromFile return err:", err)
	response, err := _CosClient.Object.Get(context.Background(), "image/test.jpg", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	//sendSms()
	//captcha()
	//sts()
	testcoswithTmpSecret()
}
