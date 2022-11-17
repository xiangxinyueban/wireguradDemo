package main

import (
	"context"
	"flag"
	"sync"
	"vpn/center"
	"vpn/daemon"
	"vpn/p2p"
)

func main() {
	mode := flag.String("mode", "center", "server mode: payload/center")

	flag.Parse()
	switch *mode {
	case "center":
		centerBootStrap()
	default:
		payloadBootStrap()
	}
}

func payloadBootStrap() {
	ctx, _ := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		p2p.BootstrapNode(ctx, "0.0.0.0", 4001)
	}()
	go func() {
		defer wg.Done()
		daemon.StartServer()
	}()
	wg.Wait()
}

func centerBootStrap() {
	center.CenterStart()
}
