package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"vpn/center/model"
	"vpn/center/serializer"
	"vpn/center/util"
)

type UserService struct {
	UserName string `form:"username" json:"username" binding:"required,min=3,max=15" example:"Anonymous"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=40" example:"Anonymous"`
	Email    string `form:"email" json:"email" binding:"required,min=7,max=40" example:"Anonymous"`
}

func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count > 0 {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "用户名已存在",
		}
	}
	user.UserName = service.UserName
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "注册失败，稍后再试",
		}
	}
	user.Email = service.Email
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "注册失败，数据库错误",
		}
	}
	return serializer.Response{
		Code:  1,
		Data:  nil,
		Msg:   "注册成功，请登录",
		Error: "",
	}
}

func (service *UserService) Login() serializer.Response {
	var user model.User
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Code:  -1,
				Error: "用户不存在",
			}
		}
		return serializer.Response{
			Code:  -1,
			Error: "登录出错",
		}
	}
	if !user.CheckPassword(service.Password) {
		return serializer.Response{
			Code:  -1,
			Error: "登录出错",
		}
	}
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		return serializer.Response{
			Code:  -1,
			Error: "Message failed",
		}
	}
	return serializer.Response{
		Code: 1,
		Data: serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:  "登录成功",
	}
}

type ActivateService struct {
	UserName     string `form:"username" json:"username" binding:"required,min=3,max=15" example:"Anonymous"`
	ActivateCode string `form:"activatecode" json:"activatecode" binding:"required,min=10,max=40 " example:"Anonymous"`
}

func (as *ActivateService) Activate() serializer.Response {
	claims, err := util.ParseActivation(as.ActivateCode)
	if err != nil || claims.Valid != "激活码" || claims.Duration <= 0 || claims.Flux <= 0 {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "激活码无效,请重新获取激活码",
		}
	}
	var order model.Order
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("user_name=?", as.UserName).First(&user).Error; err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "该激活用户不存在",
		}
	}
	order.EndTime = time.Now().Add(time.Duration(claims.Duration*24) * time.Hour)
	order.StartTime = time.Now()
	order.RemainTraffic = int64(claims.Flux) * 1024 * 1024 * 1024
	if err = model.DB.Create(&order).Error; err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "激活失败，请稍后再试",
		}
	}
	return serializer.Response{
		Code:  1,
		Data:  nil,
		Msg:   fmt.Sprintf("激活成功, 有效期: %s天; 总流量: %sG", claims.Duration, claims.Flux),
		Error: "",
	}
}

type CreateActivationService struct {
	Duration int `form:"duration" json:"duration" binding:"required" example:"Anonymous"`
	Flux     int `form:"flux" json:"flux" binding:"required " example:"Anonymous"`
}

func (as *CreateActivationService) CreateActivation() serializer.Response {
	activationCode, err := util.GenerateActivation(as.Duration, as.Flux)
	if err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: "生成激活码失败",
		}
	}

	return serializer.Response{
		Code:  1,
		Data:  nil,
		Msg:   fmt.Sprintf("激活成功, 有效期: %s天; 总流量: %sG", claims.Duration, claims.Flux),
		Error: "",
	}
}