package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"io"
)

func MD5(content []byte) []byte {
	h := md5.New()
	io.WriteString(h, string(content))
	return h.Sum(nil)
}

func HmacMD5(content []byte, key string) []byte {
	h := hmac.New(md5.New, []byte(key))
	h.Write(content)
	return h.Sum(nil)
}
