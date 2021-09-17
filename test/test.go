package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"ginfra/config"
	"ginfra/handler"
	"ginfra/models"
	"ginfra/plugin/k8sclient"
	"ginfra/tencent"
	"ginfra/utils"
	"github.com/spf13/pflag"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	kubecfg   = pflag.StringP("kubecfg", "k", "", "kubernetes apiserver config file path.")
	cmd       = pflag.StringP("command", "C", "", "tool command:createjob|getjob.")
	namespace = pflag.StringP("namespace", "n", "", "kubernetes namespace.")
)

func test_k8sclient() {
	pflag.Parse()

	kclient := &k8sclient.KubeClient{}
	if len(*kubecfg) == 0 {
		kclient.InClusterConfig()
	}

	err := kclient.WithKubeConfig(*kubecfg)
	if err != nil {
		panic(err)
	}

	if *cmd == "createjob" {
		createjob(kclient)
	} else if *cmd == "getjob" {
		getjob(kclient)
	}
}

func createjob(kclient *k8sclient.KubeClient) {
	req := &k8sclient.JobRequest{
		Namespace:               *namespace,
		JobName:                 "demo-job",
		Image:                   "xxxx",
		CpuRequest:              "700m",
		MemoryRequest:           "512Mi",
		CpuLimit:                "700m",
		MemoryLimit:             "512Mi",
		TTLSecondsAfterFinished: 300,
		Mounts: []k8sclient.VolumeMount{{
			Name:      "code",
			MountPath: "/usr/local/service/runner/code",
			HostPath:  "/tmp/test",
		}},
		Envs: map[string]string{
			"ENV_NAME": "ENV_VALUE",
		},
	}
	job, err := kclient.CreateFootballJob(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("New job: %v\n", job)
}

func getjob(kclient *k8sclient.KubeClient) {
	job, err := kclient.GetJob(context.TODO(), *namespace, "demo-job")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get job: %v\n", job)
	fmt.Printf("job status: %v\n", job.Status)
}

func testReadFile() {
	url := "http://www.golang-book.com/public/pdf/gobook.pdf"
	url = "https://test.cos.ap-shanghai.myqcloud.com/test.txt"
	cont, err := utils.ReadFile(url)
	fmt.Println("err:", err)
	fmt.Println("len(file):", len(cont))
}

func testMap() {
	tmpMap := make(map[string]interface{})
	val, ok := tmpMap["non-exist"].(uint64)
	fmt.Println(val, ok)

	val, ok = interface{}("non-number").(uint64)
	fmt.Println(val, ok)
}

func testClaims() {
	user := &models.UserAuth{
		Uid:          000001,
		IdentityType: 5,
		Identifier:   "7LSDAAEAD989QKJALJDFA",
	}

	token, err := handler.GenerateToken(user)
	fmt.Println(token, err)

	cfg, err := config.Parse("")
	if err != nil {
		panic(err)
	}

	// init
	RS256PublicKey := cfg.GetString("jwt.RS256PublicKey")
	fmt.Println(RS256PublicKey)
	claims, err := utils.ParseJWTTokenWithRS256([]byte(RS256PublicKey), token)
	data, err := handler.GetClaimData(claims)
	fmt.Println(data, err)
}

func test_jwt() {
	claims := make(map[string]interface{})
	claims["sub"] = 1
	claims["jti"] = "4a4550d0d9b3587c4f472038780452a3b17fd863c5aab7d14cca93037d49332726ab80dcbd9ddd59"
	claims["aud"] = ""
	claims["scopes"] = []interface{}{nil}
	claims["exp"] = 1630121578
	claims["iat"] = 1627529578
	claims["nbf"] = 1627529578

	publicKeyByte, _ := ioutil.ReadFile("./cert/public.key")
	privateKeyByte, _ := ioutil.ReadFile("./cert/private.key")
	token, err := utils.CreateJWTTokenFromMapWithRS256(privateKeyByte, claims)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	fmt.Println("======")

	//token = ""
	decoded, err := utils.ParseJWTTokenWithRS256(publicKeyByte, token)
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(decoded, "", "\t")
	fmt.Println(string(b))
}

func test_sts() {
	policy := "{\"statement\":[{\"action\":[\"name/cos:PutObject\",\"name/cos:PostObject\",\"name/cos:InitiateMultipartUpload\",\"name/cos:UploadPart\",\"name/cos:CompleteMultipartUpload\",\"name/cos:AbortMultipartUpload\"],\"effect\":\"allow\",\"resource\":[\"qcs::cos:ap-guangzhou:uid/APPID:bucket-name/*\"]}],\"version\":\"2.0\"}"
	resp, err := tencent.GetFederationToken("dummy", policy)
	if err != nil {
		fmt.Println(err)
		return
	}

	logstr, _ := json.Marshal(resp)
	fmt.Println(string(logstr))

	// 使用STS临时密钥
	u, _ := url.Parse("https://user-meta-1306124692.cos.ap-guangzhou.myqcloud.com")
	fmt.Println(u)
	b := &cos.BaseURL{BucketURL: u}
	_CosClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 使用GetFederationToken接口返回的临时密钥
			SecretID:     *resp.TmpSecretId,
			SecretKey:    *resp.TmpSecretKey,
			SessionToken: *resp.Token,
		},
	})

	_, err = _CosClient.Object.PutFromFile(context.Background(), "image/test.jpg", "./test.jpg", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PutFromFile succeed")
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

func test_getQQDocsToken() {
	resp, err := tencent.QQDocsToken("xxxx", "xxxx",
		"https://www.qq.com", "xxxx")
	if err != nil {
		fmt.Println(err)
		return
	}
	logstr, _ := json.Marshal(resp)
	fmt.Println(string(logstr))
}

func testPassword(pwd string) {
	err := utils.ValidatePassword(8, 32, 3, pwd)
	fmt.Println(pwd, err)
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func ObscureString(name string) string {
	runeName := []rune(name)
	if len(runeName) <= 1 {
		return name
	} else if len(runeName) == 2 {
		return "*" + string(runeName[1])
	}

	return string(runeName[0]) + "*" + string(runeName[len(runeName)-1])
}

func testMapIter() {
	type a struct {
		a string
	}

	var alist []a
	for i := 0; i < 5; i++ {
		alist = append(alist, a{
			a: strconv.Itoa(i),
		})
	}

	aMap := make(map[int]*a)
	for i, a := range alist {
		aMap[i] = &a
	}

	for k, v := range aMap {
		fmt.Println(k, v)
	}

	bMap := make(map[int]*a)
	for i := 0; i < len(alist); i++ {
		bMap[i] = &alist[i]
	}

	fmt.Println("--------------")
	for k, v := range bMap {
		fmt.Println(k, v)
	}
}

func main() {
	fmt.Println(time.Now().Format(time.RFC3339))
	// testReadFile()
	//testMap()
	//testClaims()
	//testPassword("aaaab1bB")
	//fmt.Println(VerifyEmailFormat("12345@qq.com"))  //true
	//fmt.Println(VerifyEmailFormat("12345126.@com")) //false
	//fmt.Println(ObscureString("bob"))
	//fmt.Println(ObscureString("黄生"))
	//testMapIter()
	//test_sts()
	test_getQQDocsToken()
}
