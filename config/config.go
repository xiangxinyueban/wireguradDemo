package config

import (
	"github.com/multiformats/go-multiaddr"
	"gopkg.in/ini.v1"
	"log"
	"net"
)

type LocalConfig struct {
	RPCPort           int
	ID                string
	HeartbeatInterval int
	BootstrapPeers    []multiaddr.Multiaddr
	//@TODO: RPC TLS encode?
}

type CentralConfig struct {
	Port    int
	Address net.IP
}

const CONFIG_PATH = "/var/p2pwireguard/server.ini"

var LConfig *LocalConfig

func InitLocalConfig() (cfg *LocalConfig) {
	logf, err := ini.ShadowLoad(CONFIG_PATH)
	if err != nil {
		log.Fatalln(err)
	}
	cfg.RPCPort = logf.Section("local").Key("rpc_port").MustInt(9999)
	cfg.ID = logf.Section("local").Key("id").String()
	cfg.HeartbeatInterval = logf.Section("local").Key("heartbeat_interval").MustInt(30)
	bootstrapPeers := logf.Section("local").Key("bootstrap_peer").ValueWithShadows()
	for _, v := range bootstrapPeers {
		bootstrapPeer, err := multiaddr.NewMultiaddr(v)
		if err != nil {
			log.Fatalln(err)
		}
		cfg.BootstrapPeers = append(cfg.BootstrapPeers, bootstrapPeer)
	}
	LConfig = cfg
	return
}

func InitCentralConfig() (cfg *CentralConfig) {
	logf, err := ini.ShadowLoad(CONFIG_PATH)
	if err != nil {
		log.Fatalln(err)
	}
	cfg.Port = logf.Section("center").Key("rpc_port").MustInt(9999)
	addr := logf.Section("center").Key("address").String()
	cfg.Address = net.ParseIP(addr)
	return
}
