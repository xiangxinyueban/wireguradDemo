package main

import (
	"context"
	"flag"
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
	ctx, _ := context.WithCancel(context.Background())
	p2p.BootstrapNode(ctx, "0.0.0.0", 4001)
}

func centerBootStrap() {

}
