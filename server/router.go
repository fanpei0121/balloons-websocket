package server

import (
	"os"
	"balloons/api"
	"balloons/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	// r.Use(middleware.CurrentUser())

	r.GET("test", api.Test)
	// 获取签名
	r.GET("getSign", api.GetSign)
	socket := r.Group("").Use(middleware.AuthSign())
	{
		socket.GET("/readMessage", api.ReadMessage)          // 订阅频道消息
		socket.GET("/wsPushMessage", api.WsPushMessage)      // websocket推送消息
		socket.POST("/httpPushMessage", api.HttpPushMessage) // http推送消息
	}

	return r
}
