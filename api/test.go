package api

import (
	"github.com/gin-gonic/gin"
	"balloons/service"
)

func Test(c *gin.Context)  {
	key := "86726e4356dce2d3"
	ss, err := service.GetSign(key)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	c.JSON(200, SuccessResponse(ss))
}
