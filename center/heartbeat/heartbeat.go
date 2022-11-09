package heartbeat

import (
	"context"
	"vpn/center/model"
	pb "vpn/proto"
)

type HeartBeatServer struct {
	pb.UnimplementedHeartbeatManagerServer
}

func (h *HeartBeatServer) Heartbeat(ctx context.Context, request *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	bladeInfo := request.GetBladeInfo()
	sessionInfos := request.GetSessionInfos()
	for _, session := range sessionInfos {
		var se model.Session
		model.DB.Model(&model.Session{}).Where("session_id=?", session.SessionID).First(&se)
		se.Traffic += session.TrafficUsed
		model.DB.Model(&model.Session{}).Where("session_id=?", session.SessionID).Update("traffic", se.Traffic)
	}
	var blade model.Blade
	model.DB.Model(&model.Blade{}).Where("id=?", bladeInfo.ID).First(&blade)
	model.DB.Model(&model.Blade{}).Where("id=?", bladeInfo.ID).Update("traffic", blade.Traffic)
	return &pb.HeartbeatResponse{Status: "ok"}, nil
}
