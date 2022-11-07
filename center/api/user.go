package api

import (
	"github.com/gin-gonic/gin"
	"vpn/center/service"
)

func UserRegister(c *gin.Context) {
	var userService *service.UserService
	if err := c.ShouldBind(userService); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := userService.Register()
		c.JSON(200, res)
	}
}

func UserLogin(c *gin.Context) {
	var userService *service.UserService
	if err := c.ShouldBind(userService); err != nil {
	}
}

func UserCharge(c *gin.Context) {
	var userService *service.UserService

	if err := c.ShouldBind(userService); err != nil {
	}
}
