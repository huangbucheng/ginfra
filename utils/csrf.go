package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"ginfra/log"
)

//GenCsrfToken 生成csrf token
func GenCsrfToken(key string) (string, error) {
	AesKey := []byte(key) // 对称秘钥长度必须是16的倍数
	text := fmt.Sprintf("csrf:%d", time.Now().Unix())

	encrypted, err := AesEncrypt([]byte(text), AesKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

//VerifyCsrfToken 校验csrf token
func VerifyCsrfToken(key, data string) (res bool) {
	defer func() {
		if err := recover(); err != nil {
			log.WithContext(context.Background()).Error(
				fmt.Sprintf("unexpected expection:%v", err))
			res = false
			return
		}
	}()

	AesKey := []byte(key) // 对称秘钥长度必须是16的倍数
	encrypteds, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.WithContext(context.Background()).Error(
			fmt.Sprintf("DecodeString err:%s", err.Error()))
		return
	}
	origin, err := AesDecrypt(encrypteds, AesKey)
	if err != nil {
		log.WithContext(context.Background()).Error(
			fmt.Sprintf("AesDecrypt err:%s", err.Error()))
		return
	}

	parts := strings.Split(string(origin), ":")
	if len(parts) != 2 {
		log.WithContext(context.Background()).Error(
			fmt.Sprintf("invalid AesDecrypted data:%s", origin))
		return
	}

	tokents, err := strconv.Atoi(parts[1])
	if err != nil {
		log.WithContext(context.Background()).Error(
			fmt.Sprintf("invalid AesDecrypted data:%s", origin))
		return
	}
	nowts := time.Now().Unix()

	res = math.Abs(float64(tokents)-float64(nowts)) < 600
	return
}
