package middleware

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"os"
	"balloons/api"
	"balloons/service"
)

// 验证sign是否正确
func AuthSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userParam api.UserParam
		if err := c.Bind(&userParam); err != nil {
			logs.Error("参数错误")
			c.JSON(200, api.ErrorResponse("参数错误"))
			c.Abort()
		}
		/*userModel := new (model.Users)
		user := userModel.GetUser(userParam.AppKey)*/
		secretKey := os.Getenv("SECRET_KEY")
		_, err := service.CheckSign(userParam.Sign, secretKey)
		if err != nil {
			logs.Error(err)
			c.JSON(401, api.ErrorResponse(err.Error()))
			c.Abort()
		}
		c.Next()
	}
}
