package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
	mrand "math/rand"
	"time"
	"vpn/utils/iputil"
)

var DHTBootstrapID string

// BootstrapNode relay function
func BootstrapNode(ctx context.Context, listenHost string, port int) {
	//init ipfs dht bootstrap node act as relay node. public iputil needed.

	fmt.Printf("[*] Listening on: %s with port: %d\n", listenHost, port)

	r := mrand.New(mrand.NewSource(int64(port)))

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", listenHost, port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

	_, err = dht.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", listenHost, port, host.ID())
	publicIp := iputil.GetPublicIP()
	DHTBootstrapID = fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", publicIp, port, host.ID())
	select {
	case <-ctx.Done():
		return
	default:
		log.Debug().Msg("Ping Pong DHT bootstrap running")
		time.Sleep(time.Minute)
	}
}
