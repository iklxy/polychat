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

// Register 新连接注册方法
// 如果用户已有旧连接（例如从另一个设备登录），先关闭旧连接再注册新连接
func (cm *ClientManager) Register(userID uint, conn *websocket.Conn) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()

	// 如果已存在旧连接，先关闭它，防止产生僵尸连接导致在线状态误判
	if oldConn, ok := cm.Clients[userID]; ok {
		oldConn.Close()
		fmt.Printf("用户 %d 旧连接已关闭（重新连接）\n", userID)
	}

	cm.Clients[userID] = conn
	fmt.Printf("用户已经上线 %d\n", userID)
}

// UnRegister 连接注销方法
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

// SendMessage 发送消息给指定用户
// 如果发送失败（连接已断开），自动清理该连接
func (cm *ClientManager) SendMessage(msg Message) {
	cm.Lock.RLock()
	conn, ok := cm.Clients[msg.ReceiverID]
	cm.Lock.RUnlock()

	if ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Printf("发送消息给用户 %d 失败: %v，清理断开的连接\n", msg.ReceiverID, err)
			// 写入失败说明连接已断开，清理该连接防止在线状态误报
			cm.Lock.Lock()
			// 再次检查是否是同一个连接（防止期间有新连接注册）
			if currentConn, exists := cm.Clients[msg.ReceiverID]; exists && currentConn == conn {
				conn.Close()
				delete(cm.Clients, msg.ReceiverID)
				fmt.Printf("用户 %d 的断开连接已清理\n", msg.ReceiverID)
			}
			cm.Lock.Unlock()
			return
		}
	} else {
		fmt.Printf("用户 %d 不在线\n", msg.ReceiverID)
	}
}
