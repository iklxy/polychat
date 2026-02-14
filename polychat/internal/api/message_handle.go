// Package api 提供 HTTP/WebSocket 请求处理器。
// 本文件负责处理聊天消息历史记录的 REST API 请求。
// 提供 GET /api/v1/message/history 接口，支持分页查询两个用户之间的历史消息。
package api

import (
	"net/http"
	"strconv"

	"polychat/internal/model"
	"polychat/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageHandle 消息相关的 HTTP 请求处理器。
type MessageHandle struct {
	messageService service.MessageService
}

// GetHistory 获取当前用户与指定好友之间的聊天历史记录。
//
// 请求方式: GET /api/v1/message/history
// 请求参数 (Query):
//   - target_id: 必填，聊天对象的用户ID
//   - page:      可选，页码（默认 1）
//   - page_size: 可选，每页条数（默认 50，最大 100）
//
// 响应格式:
//
//	{
//	    "code": 200,
//	    "data": {
//	        "messages": [...],      // 消息列表，按时间倒序
//	        "total": 150,           // 总消息数
//	        "page": 1,              // 当前页码
//	        "page_size": 50         // 每页条数
//	    }
//	}
//
// 错误响应:
//
//	{"code": 400, "msg": "target_id 不能为空"}
//	{"code": 500, "msg": "查询失败"}
func (h *MessageHandle) GetHistory(c *gin.Context) {
	// 从 JWT 中间件获取当前用户ID
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户未登录",
		})
		return
	}
	userID := uid.(uint)

	// 解析目标用户ID（必填参数）
	targetIDStr := c.Query("target_id")
	if targetIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "target_id 不能为空",
		})
		return
	}
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "target_id 格式错误",
		})
		return
	}

	// 解析分页参数（可选，有默认值）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	// 调用 service 层获取历史消息
	messages, total, err := h.messageService.GetHistory(userID, uint(targetID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询消息历史失败",
		})
		return
	}

	// 如果查询结果为空，返回空数组而不是 null
	if messages == nil {
		messages = []model.ChatMessage{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"messages":  messages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
