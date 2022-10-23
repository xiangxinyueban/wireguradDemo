package kernelspace

import (
	"fmt"
	"github.com/mysteriumnetwork/node/services/wireguard/wgcfg"
	"net"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestPeerStats(t *testing.T) {
	client, err := NewWireguardClient()
	if err != nil {
		fmt.Println("NewClient ERR:", err)
		return
	}
	var config wgcfg.DeviceConfig
	config.IfaceName = "wg1"
	config.Subnet = net.IPNet{IP: net.ParseIP("10.77.0.6"),
		Mask: net.IPv4Mask(255, 255, 255, 255)}
	config.Peer = wgcfg.Peer{
		PublicKey: "oW6U+PnPzyJ5MN89TiD5WGCMaQ0OwR4UqIpxJM3rbAU=",
		Endpoint: &net.UDPAddr{
			IP:   net.ParseIP("23.94.211.103"),
			Port: 41287,
		},
		AllowedIPs: []string{
			"0.0.0.0/0",
		},
		KeepAlivePeriodSeconds: 25,
	}
	config.DNS = []string{"8.8.8.8"}
	config.PrivateKey = "UHJOcIIMXFs0H2pplUdCA48YwrYbazwYODYLspMZJ20="
	config.DNSScriptDir = "/etc/openvpn/"
	defer func() {
		client.DestroyDevice(config.IfaceName)
	}()
	err = client.ConfigureDevice(config)
	if err != nil {
		fmt.Println("config ERR:", err)
	} else {
		fmt.Println(config.IfaceName, "already running")
	}
	wg.Add(1)
	go func() {

		time.AfterFunc(10*time.Minute, func() {
			stats, err := client.PeerStats("wg1")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("peer stats:", stats)
			}
			wg.Done()
		})
	}()
	wg.Wait()
}
