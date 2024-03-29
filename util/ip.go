package util

import (
	"log"
	"net"
)

func GetIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println("获取 ip 失败")
		return "0.0.0.0"
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}
