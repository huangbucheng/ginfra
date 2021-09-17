package utils

import (
	"ginfra/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

var (
	monitorSecretId  string
	monitorSecretKey string
)

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	monitorSecretId = cfg.GetString("monitor.SecretID")
	monitorSecretKey = cfg.GetString("monitor.SecretKey")
}

//PutMonitorData 上报数据到云监控
func PutMonitorData(metrics []*monitor.MetricDatum) error {
	credential := common.NewCredential(monitorSecretId, monitorSecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "monitor.tencentcloudapi.com"
	client, _ := monitor.NewClient(credential, "ap-guangzhou", cpf)

	request := monitor.NewPutMonitorDataRequest()

	request.Metrics = metrics

	_, err := client.PutMonitorData(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//fmt.Printf("An API error has returned: %s", err)
		return err
	}
	if err != nil {
		return err
	}
	//fmt.Printf("%s", response.ToJsonString())
	return nil
}
