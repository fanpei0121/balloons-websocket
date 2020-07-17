package api

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"balloons/service/impl"
)

type UserParam struct {
	AppKey       string `form:"app_key"`
	Sign         string `form:"sign"`         // 唯一标识
	ReadChannel  string `form:"readChannel"`  // 订阅读通道
	WriteChannel string `form:"writeChannel"` // 订阅写通道
	Data         string `form:"data"`         // 推送数据
}

func init() {
	go sendMessage() // 发送消息
}

func WsHandler(w http.ResponseWriter, r *http.Request) *impl.Connection {
	//	w.Write([]byte("hello"))

	var upg = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var (
		wsConn *websocket.Conn
		err    error
		conn   *impl.Connection
		// data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upg.Upgrade(w, r, nil); err != nil {
		return nil
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		conn.Close()
	}

	// 启动线程，不断发消息
	/*go func() {
		var (
			err error
		)
		for {
			if err = conn.WsConnect.WriteMessage(websocket.TextMessage, []byte("heart")); err != nil {
				conn.Close()
				break
			}
			time.Sleep(5 * time.Second)
		}
	}()*/

	return conn

}

// 消息通道
var MessageChan = make(chan UserParam, 1000)

// http推送消息
func HttpPushMessage(c *gin.Context) {
	var userParam UserParam
	if err := c.Bind(&userParam); err != nil {
		c.JSON(200, ErrorResponse("参数错误"))
		return
	}
	if userParam.AppKey == "" {
		c.JSON(200, ErrorResponse("请传app_key"))
		return
	}
	if userParam.WriteChannel == "" {
		c.JSON(200, ErrorResponse("请传推送的通道名称"))
		return
	}
	MessageChan <- userParam
	c.JSON(200, SuccessResponse("推送成功"))
}

// 推送消息
func WsPushMessage(c *gin.Context) {
	var userParam UserParam
	if err := c.Bind(&userParam); err != nil {
		c.JSON(200, ErrorResponse("参数错误"))
		return
	}
	var err error
	logs.Info(fmt.Sprintf("sign = %s channerName = %s", userParam.Sign, userParam.WriteChannel))
	conn := WsHandler(c.Writer, c.Request)
	defer func() {
		conn.Close()
	}()
	var data []byte

	for {
		if _, data, err = conn.WsConnect.ReadMessage(); err != nil {
			conn.Close()
			break
		}
		logs.Info("推送了一条数据")
		userParam.Data = string(data)
		MessageChan <- userParam
	}
}

// 发送消息
func sendMessage() {
	for {
		select {
		case userParam := <-MessageChan:
			for client := range Clients {
				if client.AppKey == userParam.AppKey && client.ReaderChannel == userParam.WriteChannel {
					if err := client.Conn.WsConnect.WriteMessage(websocket.TextMessage, []byte(userParam.Data)); err != nil {
						client.Conn.Close()
						break
					}
				}
			}
		}
	}

}

// 消息接受者
type Client struct {
	Conn          *impl.Connection // 用户websocket连接
	AppKey        string           // AppKey
	Sign          string           // 允许广播唯一用户标识
	ReaderChannel string           // 接受消息的通道
}

var Clients = make(map[Client]bool) // 用户映射

// 订阅消息
func ReadMessage(c *gin.Context) {
	var userParam UserParam
	if err := c.Bind(&userParam); err != nil {
		logs.Error("参数错误")
	}
	conn := WsHandler(c.Writer, c.Request)
	var client Client
	client.AppKey = userParam.AppKey
	client.Sign = userParam.Sign
	client.Conn = conn
	client.ReaderChannel = userParam.ReadChannel
	Clients[client] = true
	// 当函数返回时，将该用户加入退出通道，并断开用户连接
	defer func() {
		conn.Close()
		delete(Clients, client) // 将用户从映射中删除
		printClients()          // 打印Clients
	}()
	printClients() // 打印Clients
	for {
		if _, _, err := conn.WsConnect.ReadMessage(); err != nil {
			logs.Info(fmt.Sprintf("用户：%s 断开链接", userParam.AppKey))
			conn.Close()
			break
		}
	}

}

func printClients() {
	num := 1
	for key, _ := range Clients {
		fmt.Printf("用户%d %+v\n", num, key)
		num++
	}
	fmt.Println("-------")
}
