package router

import (
	"github.com/gin-gonic/gin"
	"vpn/center/api"
	"vpn/center/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	userMgmt := r.Group("user")
	{
		userMgmt.POST("register", api.UserRegister)
		userMgmt.POST("login", api.UserLogin)

		authed := userMgmt.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.POST("activate", api.Activate)
		}
	}
	r.POST("createActivation", api.CreateActivation)
	return r
}
