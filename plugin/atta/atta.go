package atta

import (
	"fmt"
	"net/url"
	"time"

	"ginfra/utils"
)

func ReportBackendRequestStatus(attaid, token string, uid, uri, code string, status, latency int) {
	_url := fmt.Sprintf(
		"https://h.trace.qq.com/kv?attaid=%s&token=%s&event_time=%s&event_code=backend_request_status"+
			"&request_path=%s&request_status=%d&request_latency=%d&uid=%s&event_result=%s",
		attaid, token, time.Now().Format(utils.TIMEFORMAT),
		url.QueryEscape(uri), status, latency, uid, code)
	utils.GetRequest(_url, nil)
	//resp, err := utils.GetRequest(_url, nil)
	//fmt.Println(err)
	//fmt.Println(string(resp))
}
