// 旧版chat的api
package api

import (
	"net/http"
	"polychat/internal/ws"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 升级器
var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectWS(c *gin.Context) {
	uid, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "用户未登录",
		})
		return
	}
	// 断言转换为uint类型
	userID := uid.(uint)
	//将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "升级为WebSocket连接失败",
		})
		return
	}
	//注册登录
	ws.ClientMgr.Register(userID, conn)

	//go协程处理连接
	go func() {
		//注销连接
		defer func() {
			ws.ClientMgr.UnRegister(userID)
		}()

		//读取信息
		for {
			msg := ws.Message{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
			//防止ID伪造
			msg.SenderID = userID
			msg.Timestamp = time.Now().Unix()
			if msg.Type == "" {
				msg.Type = ws.TypeChat
			}
			//发送信息
			ws.ClientMgr.SendMessage(msg)
		}
	}()
}
