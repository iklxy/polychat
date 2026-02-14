// Package api 提供 HTTP/WebSocket 请求处理器。
// 本文件提供带消息持久化功能的 WebSocket 连接处理器 ConnectWSWithHistory。
// 它在原有 ConnectWS 的基础上，增加了将聊天消息保存到 MongoDB 的功能，
// 使得用户可以在后续登录时查看历史消息，离线用户的消息也不会丢失。
//
// 与原有 chat.go 的区别：
//   - chat.go 中的 ConnectWS: 消息仅实时转发，不做持久化，离线消息丢失
//   - 本文件的 ConnectWSWithHistory: 每条聊天消息先持久化到 MongoDB，再转发给在线用户
package api

import (
	"fmt"
	"net/http"
	"polychat/internal/service"
	"polychat/internal/ws"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// wsUpgrader 是 WebSocket 升级器，与 chat.go 中的 upgrader 配置一致。
// 单独定义以避免修改 chat.go 中的变量。
var wsUpgrader = websocket.Upgrader{
	// 允许跨域连接
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// msgService 是消息业务服务的包级实例，供 ConnectWSWithHistory 使用。
var msgService = service.MessageService{}

// ConnectWSWithHistory 处理 WebSocket 连接请求，并在消息转发时自动持久化到 MongoDB。
// 该函数的工作流程与原有 ConnectWS 完全一致，仅在发送消息前增加了持久化步骤：
//  1. 验证用户身份（从 JWT 中间件获取 userID）
//  2. 将 HTTP 连接升级为 WebSocket 连接
//  3. 注册到全局客户端管理器
//  4. 启动 goroutine 循环读取消息
//  5. 【新增】将 chat 类型消息保存到 MongoDB
//  6. 转发消息给接收方（如果在线）
//
// 路由: GET /api/v1/chat?token=<JWT>
// 使用方式: 在 main.go 中用 ConnectWSWithHistory 替换原有的 ConnectWS 路由即可
func ConnectWSWithHistory(c *gin.Context) {
	// 从 JWT 中间件获取用户ID
	uid, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "用户未登录",
		})
		return
	}
	// 断言转换为 uint 类型
	userID := uid.(uint)

	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "升级为WebSocket连接失败",
		})
		return
	}

	// 注册到全局客户端管理器（复用已有的 ws.ClientMgr）
	ws.ClientMgr.Register(userID, conn)

	// 启动 goroutine 处理连接
	go func() {
		// 连接关闭时注销
		defer func() {
			ws.ClientMgr.UnRegister(userID)
		}()

		// 循环读取消息
		for {
			msg := ws.Message{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			// 防止 ID 伪造（与原逻辑一致）
			msg.SenderID = userID
			msg.Timestamp = time.Now().Unix()
			if msg.Type == "" {
				msg.Type = ws.TypeChat
			}

			// 【新增】将消息持久化到 MongoDB（异步，不阻塞消息转发）
			// 即使持久化失败，消息仍然会被转发给在线用户
			go func(m ws.Message) {
				if err := msgService.SaveMessage(m); err != nil {
					fmt.Printf("[MongoDB] 消息持久化失败: type=%s sender=%d receiver=%d err=%v\n",
						m.Type, m.SenderID, m.ReceiverID, err)
				} else {
					fmt.Printf("[MongoDB] 消息持久化成功: sender=%d receiver=%d\n",
						m.SenderID, m.ReceiverID)
				}
			}(msg)

			// 发送消息给接收方（复用已有的 ws.ClientMgr）
			ws.ClientMgr.SendMessage(msg)
		}
	}()
}
