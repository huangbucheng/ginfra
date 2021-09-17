package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"ginfra/utils"
)

func main() {
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
	token, err := utils.CreateJWTTokenWithRS256(privateKeyByte, claims)
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
