package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"vpn/center/model"
	"vpn/center/serializer"
	"vpn/proto"
)

type SessionService struct {
	UserName string `form:"username" json:"username" binding:"required" example:"Anonymous"`
	Country  string `form:"country" json:"country" binding:"required" example:"Anonymous"`
}

func (s *SessionService) SessionEstablish() serializer.Response {
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("user_name=?", s.UserName).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		return serializer.Response{
			Code:  -1,
			Error: "用户未注册或用户名错误",
		}
	}
	var session model.Session
	if err := model.DB.Model(&model.Session{}).Where("uid=?", user.ID).First(&session).Error; err != nil {
		//如果查询不到，返回相应的错误
		return serializer.Response{
			Code:  -1,
			Error: "该用户未激活或已过期，请获取激活码",
		}
	} else {
		var blade model.Blade
		model.DB.Model(&model.Session{}).Where(&model.Session{Status: 1, Country: s.Country}).Order("users").Find(&blade)
		address := fmt.Sprintf("%s:%v", blade.Address.String(), 8001)
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return serializer.Response{
				Code:  -1,
				Error: "服务器请求失败，请您尝试更换地点接入",
			}
		}
		defer conn.Close()
		c := pb.NewSessionManagerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if r, err := c.SessionEstablishment(ctx, &pb.SessionEstablishRequest{
			SessionID: session.SessionID,
			UserID:    uint32(session.Uid),
		}); err != nil {
			if r.GetStatus() == "OK" {
				return serializer.Response{
					Code: 1,
					Data: session.SessionID,
					Msg:  "服务器初始化成功",
				}
			}
		}
	}
	return serializer.Response{
		Code: -1,
		Msg:  "服务器初始失败，请稍后再次尝试",
	}
}

func (s *SessionService) SessionDeletion() serializer.Response {
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("user_name=?", s.UserName).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		return serializer.Response{
			Code:  -1,
			Error: "用户未注册或用户名错误",
		}
	}
	var session model.Session
	if err := model.DB.Model(&model.Session{}).Where("uid=?", user.ID).First(&session).Error; err != nil {
		//如果查询不到，返回相应的错误
		return serializer.Response{
			Code:  -1,
			Error: "该用户未激活或已过期，请获取激活码",
		}
	} else {
		var blade model.Blade
		model.DB.Model(&model.Session{}).Where(&model.Session{Status: 1, Country: s.Country}).Order("users").Find(&blade)
		address := fmt.Sprintf("%s:%v", blade.Address.String(), 8001)
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return serializer.Response{
				Code:  -1,
				Error: "服务器请求失败，请您尝试更换地点接入",
			}
		}
		defer conn.Close()
		c := pb.NewSessionManagerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if r, err := c.SessionEstablishment(ctx, &pb.SessionEstablishRequest{
			SessionID: session.SessionID,
			UserID:    uint32(session.Uid),
		}); err != nil {
			if r.GetStatus() == "OK" {
				return serializer.Response{
					Code: 1,
					Data: session.SessionID,
					Msg:  "服务器初始化成功",
				}
			}
		}
	}
	return serializer.Response{
		Code: -1,
		Msg:  "服务器初始失败，请稍后再次尝试",
	}
}
