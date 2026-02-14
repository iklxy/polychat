// 本文件负责 MongoDB 的连接初始化、集合引用管理及索引创建。
// MongoDB 用于存储用户之间的聊天历史记录，与 MySQL（存储用户和关系数据）互补。
package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB 全局变量，供其他包引用
var (
	// MongoClient 是 MongoDB 客户端实例，用于管理连接池
	MongoClient *mongo.Client

	// MongoDB 是 polychat_db 数据库实例
	MongoDB *mongo.Database

	// MongoMessageColl 是 messages 集合的引用，存储所有聊天消息
	MongoMessageColl *mongo.Collection
)

// MongoDB 连接配置常量
const (
	// mongoHost     = "47.110.94.115" // MongoDB 服务器地址（与 MySQL 同服务器）
	mongoPort     = "27017"         // MongoDB 服务端口
	mongoUser     = "admin"         // 认证用户名
	mongoPassword = "YY010303"      // 认证密码
	mongoDBName   = "polychat_db"   // 数据库名称
	mongoColl     = "messages"      // 消息集合名称
)

// InitMongoDB 初始化 MongoDB 连接并创建必要的索引。
// 该函数应在服务启动时调用，在 HTTP 服务开始监听之前完成。
// 如果连接失败或索引创建失败，程序将 panic 终止。
func InitMongoDB() {
	// 获取 MongoDB Host，默认为远程IP（用于本地开发），部署到服务器时可通过环境变量 MONGO_HOST=127.0.0.1 指定
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "47.110.94.115"
	}

	// 构建 MongoDB 连接 URI（使用认证）
	// 注意：根据您的反馈，admin 用户是在 polychat_db 数据库中验证的，
	// 所以这里 authSource 应该指向 mongoDBName (polychat_db) 而不是 admin 系统库
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s",
		mongoUser, mongoPassword, mongoHost, mongoPort, mongoDBName)

	// 设置连接超时为 10 秒，防止阻塞启动流程
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 配置客户端选项
	clientOpts := options.Client().ApplyURI(uri)

	// 建立 MongoDB 连接
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic("MongoDB 连接失败: " + err.Error())
	}

	// Ping 验证连接是否真正可用（mongo.Connect 是惰性连接，不会立即建立TCP连接）
	if err := client.Ping(ctx, nil); err != nil {
		panic("MongoDB Ping 失败（连接不可用）: " + err.Error())
	}

	// 赋值全局变量
	MongoClient = client
	MongoDB = client.Database(mongoDBName)
	MongoMessageColl = MongoDB.Collection(mongoColl)

	// 创建索引以优化查询性能
	createMessageIndexes()

	fmt.Println("MongoDB 连接成功")
}

// createMessageIndexes 为 messages 集合创建必要的索引。
// 索引策略：
//  1. 复合索引 {sender_id: 1, receiver_id: 1, timestamp: -1}
//     用于加速两个用户之间的聊天记录查询，按时间倒序排列。
//  2. 单字段索引 {timestamp: -1}
//     用于按时间排序的全局查询。
func createMessageIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			// 复合索引：优化按 sender_id + receiver_id 的查询，并按 timestamp 倒序
			Keys: bson.D{
				{Key: "sender_id", Value: 1},
				{Key: "receiver_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		{
			// 单字段索引：优化按时间排序
			Keys: bson.D{
				{Key: "timestamp", Value: -1},
			},
		},
	}

	_, err := MongoMessageColl.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		panic("MongoDB 创建索引失败: " + err.Error())
	}

	fmt.Println("MongoDB 消息索引创建成功")
}

// CloseMongoDB 优雅关闭 MongoDB 连接。
// 应在服务关闭时调用，以释放连接池资源。
func CloseMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("MongoDB 断开连接失败: %v\n", err)
		} else {
			fmt.Println("MongoDB 连接已关闭")
		}
	}
}
