package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"vpn/center/model"
	"vpn/center/serializer"
	"vpn/center/service"
)

func BladeRegister(c *gin.Context) {
	var blade service.BladeService
	if err := c.ShouldBind(&blade); err != nil {
		c.JSON(400, ErrorResponse(err))
	} else {
		res := blade.Register()
		c.JSON(200, res)
	}
}

func BladeList(c *gin.Context) {
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

	var list []model.Blade
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
