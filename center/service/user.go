package service

import (
	"github.com/jinzhu/gorm"
	"vpn/center/model"
	"vpn/center/serializer"
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
	token, err := .GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		util.LogrusObj.Info(err)
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    e.GetMsg(code),
	}
}
