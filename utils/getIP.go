package utils

import (
	"net"
)

//func GetAllIPs() []string {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		return []string{""}
//	}
//	var ipData []string
//	for _, addr := range addrs {
//		singleIP, _, err := net.ParseCIDR(addr.String())
//		if err != nil {
//			log.Fatal(err)
//		}
//		ipData = append(ipData, singleIP.String())
//	}
//	return ipData
//}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
