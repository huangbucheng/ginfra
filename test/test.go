package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"ginfra/config"
	"ginfra/handler"
	"ginfra/models"
	"ginfra/utils"
)

func testReadFile() {
	url := "http://www.golang-book.com/public/pdf/gobook.pdf"
	url = "https://ai-arena-1258274959.cos.ap-shanghai.myqcloud.com/relay_file/SH-AI-Arena-2021/20210709/35785/%E6%88%91%E7%9A%84AI_%E7%94%B5%E8%84%91AI.gz"
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
}
