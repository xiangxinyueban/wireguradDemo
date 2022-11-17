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
	//bladeMgmt := r.Group("blade")
	//{
	//	userMgmt.POST("register", api.BladeRegister)
	//	userMgmt.POST("register", api.BladeRegister)
	//	authed := userMgmt.Group("/")
	//	authed.Use(middleware.JWT())
	//	{
	//		authed.POST("activate", api.Activate)
	//	}
	//}
	r.POST("createActivation", api.CreateActivation)
	r.GET("userSum", api.UserSum)
	r.GET("bladeSum", api.BladeSum)
	r.GET("userList", api.UserList)   // pagination
	r.GET("bladeList", api.BladeList) //pagination
	r.POST("bladeRegister", api.BladeRegister)
	r.POST("sessionEstablish", api.SessionEstablish)
	//r.POST("sessionDeletion", api.SessionDeletion)
	return r
}
