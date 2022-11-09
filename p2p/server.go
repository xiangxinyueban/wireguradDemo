package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p/core/network"
	"net"
	"sync"
	"time"
	"vpn/config"
	pb "vpn/proto"
	"vpn/utils/iputil"
	"vpn/utils/key"
	"vpn/wireguard"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
)

type clientConfig struct {
	publicKey string `json:"publicKey"`
}

type serverConfig struct {
	publicKey  string `json:"publicKey"`
	assignedIP net.IP `json:"assignedIP"`
	listenPort int    `json:"listenPort"`
}

var ConfigMap map[string]*wireguard.DeviceConfig
var ServerMap map[string]*wireguard.Client
var TrafficMap map[string]uint64

var TotalBytes uint64

func init() {
	ConfigMap = make(map[string]*wireguard.DeviceConfig)
	ServerMap = make(map[string]*wireguard.Client)
	TrafficMap = make(map[string]uint64)
	TotalBytes = 0
}

func HandleSessionEstablish(sessionId string, userId string) error {
	// libp2p.New constructs a new libp2p Host. Other options can be added
	// here.
	host, err := libp2p.New(libp2p.ListenAddrs([]multiaddr.Multiaddr(nil)...))
	if err != nil {
		panic(err)
	}
	var handleStream network.StreamHandler
	done := make(chan int)
	handleStream = func(stream network.Stream) {
		log.Debug().Msg("Got a new stream!")

		// Create a buffer stream for none blocking read and write.
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		decode := json.NewDecoder(rw)
		encode := json.NewEncoder(rw)
		var clientConfig clientConfig
		err := decode.Decode(&clientConfig)
		if err != nil {
			log.Error().Msg("ClientMsg Format Error")
		}
		server, err := wireguard.NewWireguardEndpoint()
		if err != nil {
			log.Error().Err(err)
		}
		config, err := wireguard.ConfigFactory()
		if err != nil {
			log.Error().Err(err)
		}
		ConfigMap[sessionId] = config
		peerIP := &net.IPNet{
			IP:   iputil.FirstIP(config.Subnet),
			Mask: net.CIDRMask(32, 32),
		}
		config.Peer = wireguard.Peer{
			PublicKey: clientConfig.publicKey,
			AllowedIPs: []string{
				peerIP.String(),
			},
			KeepAlivePeriodSeconds: 25,
		}
		defer func() {
			delete(ConfigMap, sessionId)
			delete(ServerMap, sessionId)
			wireguard.DestroyConfig(config.IfaceName)
			server.Close()
		}()

		err = server.ConfigureDevice(*config)
		ServerMap[sessionId] = server
		if err != nil {
			fmt.Println("config ERR:", err)
		} else {
			fmt.Println(config.IfaceName, "already running")
		}
		serverPubKey, err := key.PrivateKeyToPublicKey(config.PrivateKey)
		if err != nil {
			log.Error().Err(err)
		}
		serverConfig := serverConfig{
			publicKey:  serverPubKey,
			assignedIP: peerIP.IP,
			listenPort: config.ListenPort,
		}
		encode.Encode(serverConfig)
		//go readData(rw)
		//go writeData(rw)
		rw.Flush()
		stream.Reset()
		close(done)
		// 'stream' will stay open until you close it (or the other side closes it).
	}
	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(userId), handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		return err
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return err
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	lcfg := config.InitLocalConfig()
	for _, peerAddr := range lcfg.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Debug().Err(err)
			} else {
				log.Debug().Msgf("Connection established with bootstrap node: %v", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	log.Info().Msgf("Announcing %v...", sessionId)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, sessionId)
	log.Info().Msgf("Successfully announced! %v", sessionId)

	//// Now, look for others who have announced
	//// This is like your friend telling you the location to meet you.
	//peerChan, err := routingDiscovery.FindPeers(ctx, sessionId)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for peer := range peerChan {
	//	if peer.ID == host.ID() {
	//		continue
	//	}
	//	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(userId))
	//
	//	if err != nil {
	//		continue
	//	} else {
	//		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	//		go handleStream(rw)
	//	}
	//}
	select {
	case <-done:
		break
	default:
		time.Sleep(3 * time.Second)
	}
	return nil
}

func HandleSessionDeletion(sessionID string, userID string) (uint64, error) {
	var deviceConfig *wireguard.DeviceConfig
	var ok bool
	var server *wireguard.Client
	if deviceConfig, ok = ConfigMap[sessionID]; !ok {
		return 0, nil
	}
	if server, ok = ServerMap[sessionID]; !ok {
		return 0, nil
	}
	var stats wireguard.Stats
	var err error
	stats, err = server.PeerStats("")
	if err != nil {
		return 0, err
	}
	delete(ConfigMap, sessionID)
	delete(ServerMap, sessionID)
	wireguard.DestroyConfig(deviceConfig.IfaceName)
	server.Close()
	return stats.BytesSent + stats.BytesReceived, nil
}

func PeerStats() (res []*pb.SessionInfo) {
	for sessionID, wgserver := range ServerMap {
		stats, err := wgserver.PeerStats("")
		var trafficDelta uint64
		if TrafficMap[sessionID] == 0 {
			trafficDelta = stats.BytesSent + stats.BytesReceived
		} else {
			trafficDelta = stats.BytesSent + stats.BytesReceived - TrafficMap[sessionID]
		}
		TotalBytes += trafficDelta
		TrafficMap[sessionID] = stats.BytesSent + stats.BytesReceived
		if err != nil {
			continue
		} else {
			res = append(res, &pb.SessionInfo{
				SessionID:   sessionID,
				TrafficUsed: trafficDelta,
			})
		}
	}
	return
}
