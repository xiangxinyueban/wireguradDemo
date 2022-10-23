package main

import (
	"context"
	"flag"
	"net"
	conf "vpn/config"
	"vpn/p2p"
)

func main() {
	mode := flag.String("mode", "payload", "server mode: payload/center")

	flag.Parsed()
	switch *mode {
	case "center":
		centerBootStrap()
	default:
		payloadBootStrap()
	}
}

func payloadBootStrap() {
	node := conf.Config{
		MonitorAddr: net.ParseIP(""),
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	p2p.BootstrapNode(ctx, "0.0.0.0", 4001, &node)

}

func centerBootStrap() {

}
