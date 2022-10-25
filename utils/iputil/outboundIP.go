package iputil

import (
	"github.com/pkg/errors"
	"net"
)

var checkAddress = "8.8.8.8:53"

// GetOutboundIP returns current outbound IP as string for current system
func GetOutboundIP() (string, error) {
	ip, err := getOutboundIP()
	if err != nil {
		return "", nil
	}
	return ip.String(), nil
}

func getOutboundIP() (net.IP, error) {
	ipAddress := net.ParseIP("0.0.0.0")
	localIPAddress := net.UDPAddr{IP: ipAddress}

	dialer := net.Dialer{LocalAddr: &localIPAddress}

	conn, err := dialer.Dial("udp4", checkAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine outbound IP")
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}

func GetOutboundInterface() string {
	ip, _ := getOutboundIP()
	ifaces, _ := net.Interfaces()
	for _, v := range ifaces {
		addrs, _ := v.Addrs()
		for _, addr := range addrs {
			//fmt.Println(k, "::", v, "::", addr.String(), v.Name)
			_, cidr, _ := net.ParseCIDR(addr.String())
			if cidr.Contains(ip) {
				return v.Name
			}
		}
	}
	return ""
}

func GetFreePort() (port int, err error) {
	var a *net.UDPAddr
	if a, err = net.ResolveUDPAddr("udp4", "0.0.0.0:0"); err == nil {
		var l *net.UDPConn
		if l, err = net.ListenUDP("udp4", a); err == nil {
			defer l.Close()
			return l.LocalAddr().(*net.UDPAddr).Port, nil
		}
	}
	return 0, nil
}
