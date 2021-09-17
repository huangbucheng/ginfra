package utils

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"ginfra/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	CosBucketUrl string
	CosSecretID  string
	CosSecretKey string

	_CosClient *cos.Client
	_cosonce   sync.Once
)

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	CosBucketUrl = cfg.GetString("cos.BucketURL")
	CosSecretID = cfg.GetString("cos.SecretID")
	CosSecretKey = cfg.GetString("cos.SecretKey")

}

//CosClient 获取COS Client
func CosClient() *cos.Client {
	_cosonce.Do(func() {
		u, _ := url.Parse(CosBucketUrl)
		b := &cos.BaseURL{BucketURL: u}
		_CosClient = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  CosSecretID,
				SecretKey: CosSecretKey,
			},
		})
	})

	return _CosClient
}

//PutFileToCos 上传文件
func PutFileToCos(c *cos.Client, objectKey string, localfile string) error {
	var err error
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go

	// 通过本地文件上传对象
	_, err = c.Object.PutFromFile(context.Background(), objectKey, localfile, nil)
	if err != nil {
		return err
	}
	return nil
}

//GetFileFromCos 下载文件
func GetFileFromCos(cosurl string) (string, error) {
	var name string = cosurl
	if len(CosBucketUrl) > 0 && strings.HasPrefix(cosurl, CosBucketUrl) {
		name = cosurl[len(CosBucketUrl):]
	}
	response, err := CosClient().Object.Get(context.Background(), name, nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
