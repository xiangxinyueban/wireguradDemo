package router

import (
	"github.com/gin-gonic/gin"
	"vpn/center/api"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	userMgmt := r.Group("user")
	{
		userMgmt.POST("register", api.UserRegister)
		userMgmt.POST("login", api.UserLogin)

		authed := userMgmt.Group("/")
		authed.Use()
		{
			authed.POST("charge", api.UserCharge)
		}
	}

}
