// Package model 定义应用的数据模型。
// 本文件定义 MongoDB 中存储的聊天消息文档结构。
// 消息模型与 WebSocket 协议中的 ws.Message 结构保持字段一致，
// 但增加了 MongoDB 特有的 _id 字段用于文档唯一标识。
package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// ChatMessage 表示存储在 MongoDB 中的一条聊天消息文档。
// 集合名称: messages
// 索引: {sender_id, receiver_id, timestamp} 复合索引用于加速历史查询
type ChatMessage struct {
	// ID 是 MongoDB 自动生成的文档唯一标识符 (_id)
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Type 消息类型，与 ws.TypeChat / ws.TypeHeartbeat 对应
	// 目前仅 "chat" 类型的消息会被持久化
	Type string `bson:"type" json:"type"`

	// SenderID 发送者的用户ID（对应 MySQL users 表的主键）
	SenderID uint `bson:"sender_id" json:"sender_id"`

	// ReceiverID 接收者的用户ID（对应 MySQL users 表的主键）
	ReceiverID uint `bson:"receiver_id" json:"receiver_id"`

	// Content 消息文本内容
	Content string `bson:"content" json:"content"`

	// Timestamp 消息发送时间的 Unix 时间戳（秒）
	// 由服务器在收到消息时生成，保证时间一致性
	Timestamp int64 `bson:"timestamp" json:"timestamp"`
}
