package main

import (
	"encoding/json"
	"fmt"

	"ginfra/utils"
)

func main() {
	var b []byte
	tokenresp, _ := utils.QQConnectToken("100000000", "xxxx",
		"https://www.qq.com", "xxxx")
	b, _ = json.Marshal(tokenresp)
	fmt.Println(string(b))

	openidresp, err := utils.QQConnectOpenID("xxxx")
	b, _ = json.Marshal(openidresp)
	fmt.Println(string(b), err)

	userresp, _ := utils.QQConnectUserInfo("100000000", "xxxx",
		"xxxx")
	b, _ = json.Marshal(userresp)
	fmt.Println(string(b))
}
