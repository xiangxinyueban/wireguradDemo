package api

import (
	"github.com/gin-gonic/gin"
	"vpn/center/service"
)

func Activate(c *gin.Context) {
	var activate *service.ActivateService
	if err := c.ShouldBind(activate); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := activate.Activate()
		c.JSON(200, res)
	}
}

func CreateActivation(c *gin.Context) {
	var create *service.CreateActivationService
	if err := c.ShouldBind(create); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := create.GenerateActivation()
		c.JSON(200, res)
	}
}
