// Package service 提供业务逻辑层，处于 API 处理器和 DAO 数据访问层之间。
// 本文件负责聊天消息的业务逻辑，包括消息持久化和历史记录查询。
// 业务逻辑层负责参数校验、数据转换等，将 DAO 层的原始数据操作封装为业务语义明确的方法。
package service

import (
	"fmt"

	"polychat/internal/dao"
	"polychat/internal/model"
	"polychat/internal/ws"
)

// MessageService 聊天消息业务服务。
// 提供消息持久化和历史记录查询功能。
type MessageService struct {
	messageDAO dao.MessageDAO
}

// SaveMessage 将一条 WebSocket 消息持久化到 MongoDB。
// 仅保存 type 为 "chat" 的消息，心跳等其他类型不做持久化。
//
// 参数 msg 是从 WebSocket 接收到的消息（已由服务器设置好 SenderID 和 Timestamp）。
// 返回错误信息（如果有）。
func (s *MessageService) SaveMessage(msg ws.Message) error {
	// 仅持久化聊天类型消息
	if msg.Type != ws.TypeChat {
		return nil
	}

	// 将 ws.Message 转换为 MongoDB 文档模型
	chatMsg := &model.ChatMessage{
		Type:       msg.Type,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		Timestamp:  msg.Timestamp,
	}

	if err := s.messageDAO.SaveMessage(chatMsg); err != nil {
		fmt.Printf("消息持久化失败: sender=%d, receiver=%d, err=%v\n",
			msg.SenderID, msg.ReceiverID, err)
		return err
	}

	return nil
}

// GetHistory 获取两个用户之间的聊天历史记录（分页）。
// 返回双向聊天记录（A发给B + B发给A），按时间倒序排列。
//
// 参数说明：
//   - userID:   当前登录用户的ID
//   - targetID: 聊天对象的用户ID
//   - page:     页码，从 1 开始（默认为 1）
//   - pageSize: 每页消息数量（默认为 50，最大 100）
//
// 返回值：
//   - []model.ChatMessage: 消息列表
//   - int64:               总消息数
//   - error:               错误信息
func (s *MessageService) GetHistory(userID, targetID uint, page, pageSize int) ([]model.ChatMessage, int64, error) {
	// 参数校验与默认值
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.messageDAO.GetMessageHistory(userID, targetID, page, pageSize)
}
