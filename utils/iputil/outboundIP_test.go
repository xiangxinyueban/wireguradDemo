package iputil

import (
	"fmt"
	"testing"
)

func TestGetOutboundIP(t *testing.T) {
	//outaddr, err := GetOutboundIP()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("address:", outaddr)
	//ifaces, _ := net.Interfaces()
	//for k, v := range ifaces {
	//	addrs, _ := v.Addrs()
	//	for _, addr := range addrs {
	//		fmt.Println(k, "::", v, "::", addr.String(), v.Name)
	//		_, cidr, _ := net.ParseCIDR(addr.String())
	//		if cidr.Contains(net.ParseIP(outaddr)) {
	//			fmt.Println(v.Name)
	//		}
	//	}
	//}
	iface := GetOutboundInterface()
	fmt.Println(iface)
}

func TestGetFreePort(t *testing.T) {
	port, _ := GetFreePort()
	fmt.Println(port)
}
