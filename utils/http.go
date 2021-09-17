package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
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
