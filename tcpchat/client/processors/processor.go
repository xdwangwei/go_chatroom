package processors

import (
	"awesomeProject/tcpchat/common"
	"fmt"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 处理与服务器的连接
func (this *Processor) ProcessServerConn() (err error) {
	// 创建消息接发器
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	// 循环
	for {
		// 1.接收消息,以面向对象方式调用
		msg, err := transfer.RecvDataPacket()
		if err != nil {
			return err
		}
		// 2.调用处理器处理消息
		this.dispatchProcessor(&msg)
	}
	return nil
}

// 根据不同消息类型，调用对应的处理器进行处理
func (this *Processor) dispatchProcessor(msg *common.Message) {
	switch msg.MessageType {
	// 用户状态变更类型消息 -- 用户类型处理器
	case common.NotifyUserStatusMsgType:
		up := &UserProcessor{
			Conn: this.Conn,
		}
		up.ProcessNotifyOnlineStatusResp(msg)
	// 服务器对消息的响应 -- 消息类型处理器
	case common.ChatMsgRespType:
		mp := MsgProcessor{}
		mp.ProcessMsgRespFromServer(msg)
	// 服务器转发给我的 群聊类型消息  -- 消息类型处理器
	case common.ChatMsgGroupMsgType:
		fallthrough
	// 服务器转发给我的私聊类型消息 -- 消息类型处理器
	case common.ChatMsgP2pMsgType:
		mp := MsgProcessor{}
		mp.ProcessMsgFromServer(msg)
	default:
		fmt.Println("未知的消息体类型,无法处理!", msg)
	}
}
