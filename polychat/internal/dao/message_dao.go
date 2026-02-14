// Package dao 提供数据访问层，封装数据库操作。
// 本文件负责 MongoDB 中聊天消息的增删查操作。
// 所有方法均操作 database.MongoMessageColl 集合。
package dao

import (
	"context"
	"time"

	"polychat/internal/model"
	"polychat/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageDAO 聊天消息数据访问对象。
// 封装了对 MongoDB messages 集合的所有操作。
type MessageDAO struct{}

// SaveMessage 将一条聊天消息保存到 MongoDB。
// 参数 msg 是要保存的消息，ID 字段会由 MongoDB 自动生成。
// 返回错误信息（如果有）。
func (d *MessageDAO) SaveMessage(msg *model.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.MongoMessageColl.InsertOne(ctx, msg)
	return err
}

// GetMessageHistory 获取两个用户之间的聊天历史记录（分页）。
// 查询逻辑：查找 (sender=userA AND receiver=userB) OR (sender=userB AND receiver=userA) 的所有消息，
// 即双向聊天记录。结果按时间戳倒序排列（最新的在前）。
//
// 参数说明：
//   - userID:   当前用户ID
//   - targetID: 聊天对象ID
//   - page:     页码，从 1 开始
//   - pageSize: 每页消息数量
//
// 返回值：
//   - []model.ChatMessage: 消息列表（按时间倒序）
//   - int64:               总消息数（用于前端分页计算）
//   - error:               错误信息
func (d *MessageDAO) GetMessageHistory(userID, targetID uint, page, pageSize int) ([]model.ChatMessage, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件：双向匹配（A发给B 或 B发给A）
	filter := bson.M{
		"$or": bson.A{
			bson.M{"sender_id": userID, "receiver_id": targetID},
			bson.M{"sender_id": targetID, "receiver_id": userID},
		},
	}

	// 查询总数（用于分页）
	total, err := database.MongoMessageColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 计算跳过的文档数
	skip := int64((page - 1) * pageSize)

	// 设置查询选项：按时间倒序，分页
	findOpts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	// 执行查询
	cursor, err := database.MongoMessageColl.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 解码结果
	var messages []model.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
