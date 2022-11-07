package config

import (
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"net"
)

type Config struct {
	MonitorAddr        net.IP
	DHTBootStrapPeerID string
}

type P2pServerConfig struct {
	BootstrapPeers []multiaddr.Multiaddr
}

var P2pConfig P2pServerConfig

type RpcConfigStruc struct {
	Port int
}

var RpcConfig RpcConfigStruc

func init() {
	cfg, err := ini.ShadowLoad("/var/vpn/config.ini")
	if err != nil {
		log.Fatal().Msg("config not exist please touch /var/config/p2p.ini first ")
	}
	bootstrappeers := cfg.Section("p2p").Key("bootstrappeers").ValueWithShadows()
	for _, v := range bootstrappeers {
		addr, err := multiaddr.NewMultiaddr(v)
		if err != nil {
			log.Fatal().Msg("bootstrappers multiaddress format invalid")
		}
		P2pConfig.BootstrapPeers = append(P2pConfig.BootstrapPeers, addr)
	}
	RpcConfig.Port = cfg.Section("rpc").Key("port").MustInt(8888)
}
