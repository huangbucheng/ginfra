package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Success is indicated with 2xx status codes:
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		// You may read / inspect response body
		return nil, fmt.Errorf("HTTP Response Status:%d, Body:%s",
			response.StatusCode, body)
	}

	return body, nil
}

func PostRequest(url string, headers map[string]string, reqbody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbody))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Success is indicated with 2xx status codes:
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		// You may read / inspect response body
		return nil, fmt.Errorf("HTTP Response Status:%d, Body:%s",
			response.StatusCode, body)
	}

	return body, nil
}
