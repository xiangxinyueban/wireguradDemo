package service

import (
	"fmt"
	"strconv"
	"time"
	"vpn/center/lcache"
	"vpn/center/model"
	"vpn/center/serializer"
	"vpn/center/snow"
	"vpn/center/util"
)

type ActivateService struct {
	UserName string `form:"username" json:"username" binding:"required,min=3,max=15" example:"Anonymous"`
	Code     string `form:"code" json:"code" binding:"required,min=8,max=40" example:"Anonymous"`
}

func (as *ActivateService) Activate() serializer.Response {
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("user_name=?", as.UserName).First(&user).Count(&count)
	if count == 0 {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "非法参数,待激活用户不存在",
		}
	}
	token, exist := lcache.GetCache(as.Code)
	if !exist {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "激活码不存在或者超时",
		}
	}
	claims, err := util.ParseActivation(token.(string))
	if err != nil || claims.Type != "激活码" || claims.Duration <= 0 || claims.Flux <= 0 {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "激活码有误",
		}
	}
	var order model.Order
	order.RemainTraffic = int64(claims.Flux) * 1024 * 1024 * 1024
	order.StartTime = time.Now()
	order.EndTime = order.StartTime.Add(time.Duration(claims.Duration) * 24 * time.Hour)
	order.Uid = user.UserId
	if err := model.DB.Create(&order).Error; err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "激活失败，数据库问题",
		}
	}
	return serializer.Response{
		Code:  1,
		Data:  nil,
		Msg:   fmt.Sprintf("激活成功, 有效时长: %v, 总流量: %v", claims.Duration, claims.Flux),
		Error: "",
	}
}

type CreateActivationService struct {
	Duration int `form:"duration" json:"duration" binding:"required" example:"Anonymous"`
	Flux     int `form:"flux" json:"flux" binding:"required" example:"Anonymous"`
}

func (cas *CreateActivationService) GenerateActivation() serializer.Response {
	activateToken, err := util.GenerateActivation(cas.Duration, cas.Flux)
	if err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "生成激活码失败",
		}
	}
	activationCode := strconv.FormatInt(snow.SN.GetID(), 10)
	lcache.AddCache(activationCode, activateToken, 10*time.Minute)
	return serializer.Response{
		Code: 1,
		Msg:  activationCode,
	}
}
