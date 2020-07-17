package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"balloons/service"
)

// 获取签名
func GetSign(c *gin.Context) {
	secretKey := os.Getenv("SECRET_KEY")
	sign, err := service.GetSign(secretKey)
	if err != nil {
		c.JSON(200, ErrorResponse("获取Sign失败"))
		return
	}
	res := make(map[string]string)
	res["sign"] = sign
	c.JSON(200, SuccessResponse(Success, res))
}
