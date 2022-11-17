package center

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"vpn/center/heartbeat"
	"vpn/center/model"
	"vpn/center/router"
	pb "vpn/proto"
)

func CenterStart() {
	model.InitDB()
	r := router.NewRouter()
	go func() {
		ls, err := net.Listen("tcp4", "0.0.0.0:9001")
		if err != nil {
			log.Fatal().Msgf("rpc server start failed: %v", err)
		}

		s := grpc.NewServer()
		log.Debug().Msg("heartbeat RPC server ready")
		pb.RegisterHeartbeatManagerServer(s, &heartbeat.HeartBeatServer{})
		if err := s.Serve(ls); err != nil {
			log.Fatal().Msgf("rpc server start failed: %v", err)
		}
	}()
	r.Run(":9090")
}
