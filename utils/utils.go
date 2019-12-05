package utils

import (
	"net"
	"strings"
)

func GetLocalAddress() (addr []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	addr = make([]string, 0)
	for _, lip := range addrs {

		ipstr := lip.String()
		if strings.Index(ipstr, "127.0.0.1") >= 0 {
			continue
		}
		if strings.Index(ipstr, "::1") >= 0 {
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
