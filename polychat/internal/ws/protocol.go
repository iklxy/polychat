package ws

// 定义前后端通信的消息结构体

// 消息类型
const (
	TypeChat          = "chat"           // 聊天消息
	TypeHeartbeat     = "heartbeat"      // 心跳消息
	TypeFriendRequest = "friend_request" // 好友请求通知
	TypeFriendAccept  = "friend_accept"  // 好友接受通知
)

type Message struct {
	Type       string `json:"type"`        //消息类型
	SenderID   uint   `json:"sender_id"`   //发送者ID
	ReceiverID uint   `json:"receiver_id"` //接收者ID
	Content    string `json:"content"`     //消息内容
	Timestamp  int64  `json:"timestamp"`   //消息时间戳
}
