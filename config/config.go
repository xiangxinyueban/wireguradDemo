package config

import "net"

type Config struct {
	MonitorAddr        net.IP
	DHTBootStrapPeerID string
}
