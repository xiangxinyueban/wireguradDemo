package main

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"vpn/center/heartbeat"
	"vpn/center/model"
	"vpn/center/router"
	pb "vpn/proto"
)

// register -> login -> combo list -> save combo type in user config
// combo list -> create order -> order success ->
// server list(country) -> establish session -> Success page()
// delete session

// mysql struct
// user:
// UserName, Password, UserId, Email
// task:
// ComboType, RemainTraffic, ExpireTime, UserId,
// role?
//

func main() {
	model.InitDB()
	r := router.NewRouter()
	go func() {
		ls, err := net.Listen("tcp4", "0.0.0.0:9191")
		if err != nil {
			log.Fatal().Msgf("rpc server start failed: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterHeartbeatManagerServer(s, &heartbeat.HeartBeatServer{})
		if err := s.Serve(ls); err != nil {
			log.Fatal().Msgf("rpc server start failed: %v", err)
		}
	}()
	r.Run(":9090")
}
