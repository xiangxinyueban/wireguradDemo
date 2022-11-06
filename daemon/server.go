package daemon

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
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
	listen, err := net.Listen("tcp", fmt.Sprintf("%v", config.RpcConfig.Port))
	if err != nil {
		log.Fatal().Msg("RPC listen failed")
	}
	s := grpc.NewServer()
	pb.RegisterSessionManagerServer(s, &SessionManager{})
	log.Info().Msg("Serving gRPC on 0.0.0.0" + string(config.RpcConfig.Port))
	if err := s.Serve(listen); err != nil {
		log.Fatal().Err(err).Msg("failed to serve RPC")
	}
}
