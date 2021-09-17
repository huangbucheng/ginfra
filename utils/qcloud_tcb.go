package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ginfra/config"
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

var (
	signhost  string = "api.tcloudbase.com"
	algorithm string = "TC3-HMAC-SHA256"
	service   string = "tcb"
	//version   string = "2017-03-12"
	//region    string = "ap-shanghai"

	host      string = "tcb-api.tencentcloudapi.com"
	secretId  string
	secretKey string
	//envId      string //= "ai-arena-4gfpsmhn2217a420"
	//collection string //= "shanghai-2021-live-bullet"
)

func init() {
	var err error
	var cfg *config.Config
	cfg, err = config.Parse("")
	if err != nil {
		panic(err)
	}

	//envId = cfg.GetString("tcb.envId")
	//collection = cfg.GetString("tcb.collection")
	secretId = cfg.GetString("tcb.secretId")
	secretKey = cfg.GetString("tcb.secretKey")
}

func signature(secretId string, secretKey string, timestamp int64) string {
	// fmt.Println("~~~~step 1: build canonical request string")
	httpRequestMethod := "POST"
	canonicalURI := fmt.Sprintf("//%s/", signhost)
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json; charset=utf-8\n" + "host:" + signhost + "\n"
	signedHeaders := "content-type;host"
	hashedRequestPayload := sha256hex("")
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload)
	// fmt.Println(canonicalRequest)

	// fmt.Println("~~~~step 2: build string to sign")
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		timestamp,
		credentialScope,
		hashedCanonicalRequest)
	// fmt.Println(string2sign)

	// fmt.Println("~~~~step 3: sign string")
	secretDate := hmacsha256(date, "TC3"+secretKey)
	secretService := hmacsha256(service, secretDate)
	secretSigning := hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacsha256(string2sign, secretSigning)))
	// fmt.Println(signature)

	// fmt.Println("~~~~step 4: build authorization")
	authorization := fmt.Sprintf("1.0 %s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		secretId,
		credentialScope,
		signedHeaders,
		signature)
	// fmt.Println(authorization)
	return authorization
}

func tcbget(url, authorization string, timestamp int64) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("X-CloudBase-TimeStamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-CloudBase-Authorization", authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respByte))
	return respByte, nil
}

func tcbpost(url, authorization string, timestamp int64, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("X-CloudBase-TimeStamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-CloudBase-Authorization", authorization)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respByte))
	return respByte, nil
}

type tcbInsertDocRequest struct {
	Data []string `json:"data"`
}

type tcbQueryDocResponse struct {
	RequestId string          `json:"requestId"`
	Data      tcbQueryDocData `json:"data"`
}
type tcbQueryDocData struct {
	Offset int      `json:"offset"`
	Limit  int      `json:"limit"`
	List   []string `json:"list"`
}

type tcbInsertDocResponse struct {
	RequestId string           `json:"requestId"`
	Data      tcbInsertDocData `json:"data"`
}
type tcbInsertDocData struct {
	InsertedIds []string `json:"insertedIds"`
}

//QueryDocument 查询腾讯云云开发数据库
func QueryDocument(envId, collection, docId string) ([]string, error) {
	var timestamp int64 = time.Now().Unix()
	var authorization string = signature(secretId, secretKey, timestamp)

	url := fmt.Sprintf("https://%s/api/v2/envs/%s/databases/%s/documents/%s",
		host, envId, collection, docId)
	resp, err := tcbget(url, authorization, timestamp)
	if err != nil {
		return []string{}, err
	}

	var response tcbQueryDocResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []string{string(resp)}, err
	}
	return response.Data.List, nil
}

//InsertDocuments 往腾讯云云开发数据库插入数据
func InsertDocuments(envId, collection string, docs [][]byte) ([]string, error) {
	var timestamp int64 = time.Now().Unix()
	var authorization string = signature(secretId, secretKey, timestamp)

	body := &tcbInsertDocRequest{
		Data: []string{},
	}
	for _, doc := range docs {
		body.Data = append(body.Data, string(doc))
	}
	b, _ := json.Marshal(body)

	url := fmt.Sprintf("https://%s/api/v2/envs/%s/databases/%s/documents",
		host, envId, collection)
	resp, err := tcbpost(url, authorization, timestamp, b)
	if err != nil {
		return []string{}, err
	}

	var response tcbInsertDocResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []string{}, err
	}
	if len(response.Data.InsertedIds) <= 0 {
		return []string{}, fmt.Errorf("tcb return no insertedIds:%v", response)
	}
	return response.Data.InsertedIds, nil
}
