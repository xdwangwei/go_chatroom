package common

// 消息类型定义
const (
	LoginMsgType            = "loginMessage"
	LoginRespMsgType        = "loginRespMessage"
	RegisterMsgType         = "registerMessage"
	RegisterRespMsgType     = "registerRespMessage"
	NotifyUserStatusMsgType = "notifyUserStatusMsgType"
	ChatMsgGroupMsgType     = "GroupMessage"
	ChatMsgP2pMsgType       = "P2pMessage"
	ChatMsgRespType         = "chatMsgRespType"
)

// 用户在线状态定义
const (
	UserStatusOnline = iota
	UserStatusOffline
)

// 客户端和服务器传输消息体
type Message struct {
	MessageType string
	Data        string
}

// 用户登录所用的数据对象
type LoginBody struct {
	UserId   int
	Password string
}

// 服务器响应登录请求返回的数据对象
type LoginResp struct {
	Code int
	Msg  string
	// 当前在线列表
	OnlineList []int
}

// 用户注册所用的数据对象
type RegisterBody struct {
	UserId   int
	Password string
	Username string
}

// 服务器响应注册请求返回的数据对象
type RegisterResp struct {
	Code int
	Msg  string
}

// 服务器通知用户在线状态消息
type UserOnlineStatusBody struct {
	UserId int
	Status int
}

// 消息体

type ChatMsgBody struct {
	// 发送者id,就是自己
	SenderId int
	// 接收者id,如果是群聊消息,就不用管这个字段,消息体最终要封装成Message类型,通过type字段区分
	ReceiverId int
	// 消息内容
	Content string
}

// 消息响应体
type ChatMsgResp struct {
	Code int
	Msg  string
	// 如果发送失败,则把原消息体返回,实际应该是 ChatMsgBody 类型,string是为了传输方便
	Data string
}
