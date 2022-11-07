package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"vpn/center/util"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errMsg string
		var data interface{}
		token := c.GetHeader("Authorization")
		if token == "" {

		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				errMsg = "authentication failed"
			} else if time.Now().Unix() > claims.ExpiresAt {
				errMsg = "authentication timeout"
			}
		}
		if errMsg != "" {
			c.JSON(400, gin.H{
				"code":  -1,
				"msg":   "",
				"data":  data,
				"error": errMsg,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
