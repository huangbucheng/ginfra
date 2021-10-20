package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"
)

const TIMEFORMAT = "2006-01-02 15:04:05"

//GetLocalAddress 获取本地IP
func GetLocalAddress() (addr []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	addr = make([]string, 0)
	for _, lip := range addrs {

		ipstr := lip.String()
		if strings.Contains(ipstr, "127.0.0.1") {
			continue
		}
		if strings.Contains(ipstr, "::1") {
			continue
		}
		index := strings.Index(ipstr, "/")
		ip := ipstr
		if index > 0 {
			ip = ipstr[0:index]
		}
		if strings.ContainsRune(ip, rune(':')) {
			continue
		}
		addr = append(addr, ip)
	}
	return addr, nil
}

//GetEth0Addr 获取eth0网卡IP
func GetEth0Addr() string {
	var address string
	addrs, err := GetLocalAddress()
	if err != nil || len(addrs) == 0 {
		address = "InvalidAddr"
	} else {
		address = addrs[0]
	}
	return address
}

//给字符串生成md5
//@params str 需要加密的字符串
//@params salt interface{} 加密的盐
//@return str 返回md5码
func Md5Crypt(str string, salt ...interface{}) string {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(strings.Join(slice, "%v")+str, salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

//GeneratePassword 生成随机密码
func GeneratePassword(length int) string {
	if length < 8 {
		length = 8
	} else if length > 32 {
		length = 32
	}

	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~!@#$?"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits

	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

// t, err = ParseWithLocation("Asia/Shanghai", "2020-07-29 15:04:05")
func ParseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		println(err.Error())
		return time.Time{}, err
	} else {
		lt, _ := time.ParseInLocation(TIMEFORMAT, timeStr, l)
		fmt.Println(locationName, lt)
		return lt, nil
	}
}

//StringInSlice 判断字符串是否在slice中
func StringInSlice(element string, elements []string) (isIn bool) {
	for _, item := range elements {
		if element == item {
			isIn = true
			return
		}
	}
	return
}

//IntInSlice 判断数字是否在slice中
func IntInSlice(element int, elements []int) (isIn bool) {
	for _, item := range elements {
		if element == item {
			isIn = true
			return
		}
	}
	return
}

//UnmarshalBigInt due to json.Unmarshal int64 to json float64 would lose precision
func DecodeJsonWithInt64(jsonBytes []byte, result interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()
	return decoder.Decode(&result)
}

const (
	levelD = iota
	LevelC
	LevelB
	LevelA
	LevelS
)

func ValidatePassword(minLength, maxLength, minLevel int, pwd string) error {
	if len(pwd) < minLength {
		return fmt.Errorf("BAD PASSWORD: The password is shorter than %d characters", minLength)
	}
	if len(pwd) > maxLength {
		return fmt.Errorf("BAD PASSWORD: The password is logner than %d characters", maxLength)
	}

	var level int = levelD
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}

	if level < minLevel {
		return fmt.Errorf("The password does not satisfy the current policy requirements. ")
	}
	return nil
}

//ObscureName 隐藏名称部分字符
func ObscureName(name string) string {
	runeName := []rune(name)
	if len(runeName) <= 1 {
		return name
	} else if len(runeName) == 2 {
		return "*" + string(runeName[1])
	}

	return string(runeName[0]) + "*" + string(runeName[len(runeName)-1])
}

//TruncateString 截取一定长度字符串
func TruncateString(raw string, length int) string {
	runes := []rune(raw)
	runeslen := len(runes)
	return string(runes[:Min(length, runeslen)])
}
