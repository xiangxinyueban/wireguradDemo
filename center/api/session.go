package api

import (
	"github.com/gin-gonic/gin"
	"vpn/center/service"
)

func SessionEstablish(c *gin.Context) {
	var session service.SessionService

	if err := c.ShouldBind(&session); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := session.SessionEstablish()
		c.JSON(200, res)
	}
}
