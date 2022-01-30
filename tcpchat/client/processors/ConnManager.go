package processors

import "net"

// 管理和服务器的通信涉及到的数据等，目前 主要是在线用户列表

// 提供全局对象
var (
	GlobalConnManager *ConnManager
)

type ConnManager struct {
	// 在线用户列表,只保存在线的
	userList []int
	// 当前用户的信息,当前用户和服务器的连接,全局保存后,便于在其他地方获取
	UserId int
	Conn   net.Conn
}

// 自动被调用
func init() {
	if GlobalConnManager == nil {
		GlobalConnManager = &ConnManager{
			userList: make([]int, 0, 10),
		}
	}
}

// 加入某个用户
func (this *ConnManager) AddUser(userId int) {
	this.userList = append(this.userList, userId)
}

// 加入多个个用户
func (this *ConnManager) AddUserBatch(userIdList []int) {
	this.userList = append(this.userList, userIdList...)
}

// 移除某个用户
func (this *ConnManager) DelUser(userId int) {
	if len(this.userList) == 0 {
		return
	}
	index := 0
	for idx, uId := range this.userList {
		if uId == userId {
			index = idx
			break
		}
	}
	// 如果是最后一个
	if index == len(this.userList)-1 {
		this.userList = this.userList[:index]
		// 如果是中间部分某个元素,保证index+1不会引起越界
	} else {
		this.userList = append(this.userList[:index], this.userList[index+1])
	}
}

// 获取全部
func (this *ConnManager) GetAllUser() []int {
	return this.userList
}
