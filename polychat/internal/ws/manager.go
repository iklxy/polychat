package ws

//管理所有连接
import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// ClientManager 客户端管理器
type ClientManager struct {
	Clients map[uint]*websocket.Conn //存储所有的连接，key为UserID，value为连接
	Lock    sync.RWMutex             // 读写锁，保护Clients map的并发访问
}

// 全局唯一的客户端管理器实例
var ClientMgr = &ClientManager{
	Clients: make(map[uint]*websocket.Conn),
}

// 新连接注册方法
func (cm *ClientManager) Register(userID uint, conn *websocket.Conn) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()

	cm.Clients[userID] = conn
	fmt.Printf("用户已经上线 %d\n", userID)
}

// 连接注销方法
func (cm *ClientManager) UnRegister(userID uint) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()

	if conn, ok := cm.Clients[userID]; ok {
		conn.Close()
		delete(cm.Clients, userID)
		fmt.Printf("用户已经下线 %d\n", userID)
	}
}

// IsUserOnline 检查用户是否在线
func (cm *ClientManager) IsUserOnline(userID uint) bool {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	_, ok := cm.Clients[userID]
	return ok
}

// 发送信息方法
func (cm *ClientManager) SendMessage(msg Message) {
	cm.Lock.RLock()
	conn, ok := cm.Clients[msg.ReceiverID]
	cm.Lock.RUnlock()

	if ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			//发送失败
			fmt.Printf("发送消息给用户 %d 失败: %v\n", msg.ReceiverID, err)
			return
		}
	} else {
		//用户不在线
		fmt.Printf("用户 %d 不在线\n", msg.ReceiverID)
	}
}
