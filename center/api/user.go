package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
	//var userService service.UserService
	//if err := c.ShouldBind(&userService); err != nil {
	//	c.JSON(400, ErrorResponse(err))
	//} else {
	res := bladeSum()
	c.JSON(200, res)
	//}
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

func UserList(c *gin.Context) {
	//分页查询
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	//判断是否需要分页
	if limit == 0 {
		limit = -1
	}

	if page == 0 {
		page = -1
	}

	offsetVal := (page - 1) * limit
	if page == -1 && limit == -1 {
		offsetVal = -1
	}

	var list []model.User
	//返回总数
	//查询数据库
	var total int64
	model.DB.Model(list).Count(&total).Limit(limit).Offset(offsetVal).Find(&list)
	c.JSON(200, serializer.Response{
		Code: 1,
		Data: gin.H{
			"list":     list,
			"total":    total,
			"pageNum":  page,
			"pageSize": limit,
		},
		Msg:   "查询成功",
		Error: "",
	})
}

//func userList() serializer.Response {
//
//}
