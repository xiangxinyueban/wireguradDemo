package api

import (
	"github.com/gin-gonic/gin"
	"vpn/center/model"
	"vpn/center/serializer"
	"vpn/center/service"
)

func UserRegister(c *gin.Context) {
	var userService service.UserService
	if err := c.ShouldBind(&userService); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := userService.Register()
		c.JSON(200, res)
	}
}

func UserLogin(c *gin.Context) {
	var userService service.UserService
	if err := c.ShouldBind(&userService); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := userService.Login()
		c.JSON(200, res)
	}
}

func UserSum(c *gin.Context) {
	//@TODO: Permission management
	//var userService service.UserService
	//if err := c.ShouldBind(&userService); err != nil {
	//	c.JSON(400, ErrorResponse(err))
	//} else {
	res := userSum()
	c.JSON(200, res)
	//}
}

func BladeSum(c *gin.Context) {
	var userService service.UserService
	if err := c.ShouldBind(&userService); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := bladeSum()
		c.JSON(200, res)
	}
}

func userSum() serializer.Response {
	var activeCount, inactiveCount int64
	model.DB.Model(&model.User{}).Where("status=?", 1).Count(&activeCount)
	model.DB.Model(&model.User{}).Where("status=?", 0).Count(&inactiveCount)
	return serializer.Response{
		Code: 1,
		Data: serializer.Sum{Active: activeCount, InActive: inactiveCount},
		Msg:  "success",
	}
}

func bladeSum() serializer.Response {
	var activeCount, inactiveCount int64
	model.DB.Model(&model.Blade{}).Where("status=?", 1).Count(&activeCount)
	model.DB.Model(&model.Blade{}).Where("status=?", 0).Count(&inactiveCount)
	return serializer.Response{
		Code: 1,
		Data: serializer.Sum{Active: activeCount, InActive: inactiveCount},
		Msg:  "success",
	}
}
