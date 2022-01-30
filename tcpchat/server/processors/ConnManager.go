package processors

import "net"

// 全局变量，全局只用一个这个对象
var (
	GlobalConnManager *ConnManager
)

// 结构体，用于管理所有已连接的客户端
type ConnManager struct {
	clientMap map[int]*UserProcessor
}

// 这个函数会在此包被引入时自动调用
func init() {
	// 初始化这个全局变量
	if GlobalConnManager == nil {
		GlobalConnManager = &ConnManager{
			// 设定初始容量为1024
			clientMap: make(map[int]*UserProcessor, 1024),
		}
	}
}

// 添加一个新用户
func (this *ConnManager) AddNewUser(userId int, up *UserProcessor) {
	this.clientMap[userId] = up
}

// 获取一个用户的处理器
func (this *ConnManager) GetProcessorByUserId(userId int) (up *UserProcessor) {
	return this.clientMap[userId]
}

// 获取一个用户的处理器
func (this *ConnManager) GetProcessorByConn(conn net.Conn) (up *UserProcessor) {
	for _, up := range this.clientMap {
		if up.Conn == conn {
			return up
		}
	}
	return nil
}

// 删除一个客户端连接
func (this *ConnManager) DelUserByUid(userId int) {
	delete(this.clientMap, userId)
}

// 删除一个客户端连接
func (this *ConnManager) DelUserByConn(conn net.Conn) {
	key := -1
	for id, up := range this.clientMap {
		if up.Conn == conn {
			key = id
			break
		}
	}
	delete(this.clientMap, key)
}

// 获取全部用户Id
func (this *ConnManager) GetAllUser() (userSlice []int) {
	for id, _ := range this.clientMap {
		userSlice = append(userSlice, id)
	}
	return
}
