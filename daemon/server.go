package daemon

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
	"vpn/config"
	"vpn/p2p"
	"vpn/proto"
)

type SessionManager struct {
	pb.UnimplementedSessionManagerServer
}

func (se *SessionManager) SessionEstablishment(ctx context.Context, in *pb.SessionEstablishRequest) (*pb.SessionEstablishResponse, error) {
	sessionID := in.GetSessionID()
	userID := in.GetUserID()
	if err := p2p.HandleSessionEstablish(sessionID, userID); err != nil {
		return &pb.SessionEstablishResponse{Status: "Session Establishment Failed"}, err
	}
	return &pb.SessionEstablishResponse{Status: "Session Establishment Success"}, nil
}

func (se *SessionManager) SessionDeletion(ctx context.Context, in *pb.SessionDeletionRequest) (*pb.SessionDeletionResponse, error) {
	sessionID := in.GetSessionID()
	userID := in.GetUserID()
	bytesUsed, err := p2p.HandleSessionDeletion(sessionID, userID)
	if err != nil {
		return &pb.SessionDeletionResponse{Status: bytesUsed}, err
	}
	return &pb.SessionDeletionResponse{Status: bytesUsed}, nil
}

func StartServer() {
	lcfg := config.InitLocalConfig()
	listen, err := net.Listen("tcp", fmt.Sprintf("%v", lcfg.RPCPort))
	if err != nil {
		log.Fatal().Msg("RPC listen failed")
	}
	s := grpc.NewServer()
	pb.RegisterSessionManagerServer(s, &SessionManager{})
	done := make(chan int)
	defer func() {
		close(done)
	}()
	go func() {
		ticker := time.NewTicker(time.Duration(lcfg.HeartbeatInterval) * time.Second)
		select {
		case <-ticker.C:
			log.Debug().Msg("heartbeat start, interval: " + string(lcfg.HeartbeatInterval) + "s")
			heartbeat()
		case <-done:
			log.Debug().Msg("heartbeat end")
			break
		}
	}()
	log.Info().Msg("Serving gRPC on 0.0.0.0" + string(lcfg.RPCPort))
	if err := s.Serve(listen); err != nil {
		log.Fatal().Err(err).Msg("failed to serve RPC")
	}
}

func heartbeat() {
	centralCfg := config.InitCentralConfig()
	address := centralCfg.Address.String() + strconv.Itoa(centralCfg.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err)
	}
	defer conn.Close()
	c := pb.NewHeartbeatManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sessionInfos := p2p.PeerStats()
	bladeInfo := &pb.BladeInfo{
		ID:          config.LConfig.ID,
		TrafficUsed: p2p.TotalBytes,
		BootstrapID: p2p.DHTBootstrapID, //local blade as DHT bootstrap node.
	}
	r, err := c.Heartbeat(ctx, &pb.HeartbeatRequest{
		BladeInfo:    bladeInfo,
		SessionInfos: sessionInfos,
	})
	if err != nil {
		log.Error().Msg("heartbeat send failed; please check center status")
	}
	log.Printf("HeartBeat OK: %s", r.GetStatus())
}
