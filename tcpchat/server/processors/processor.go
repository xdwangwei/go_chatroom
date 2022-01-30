package processors

import (
	"awesomeProject/tcpchat/common"
	"fmt"
	"net"
)

// 将方法与结构体进行绑定,以面向对象方式调用
type Processor struct {
	Conn net.Conn
}

// 处理与客户端的连接
func (this *Processor) ProcessClientConn() (err error) {
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
	// 用户登录请求 -- 用户类型处理器
	case common.LoginMsgType:
		userProcessor := &UserProcessor{
			Conn: this.Conn,
		}
		userProcessor.ProcessLoginReq(msg)
	// 用户注册请求 -- 用户类型处理器
	case common.RegisterMsgType:
		// 创建对应类型处理器对象,调用其方法
		userProcessor := &UserProcessor{
			Conn: this.Conn,
		}
		userProcessor.ProcessRegisterReq(msg)
	// 群聊类型消息 -- 消息类型处理器
	case common.ChatMsgGroupMsgType:
		msgProcessor := MsgProcessor{
			Conn: this.Conn,
		}
		msgProcessor.processGroupMsg(msg)
	// 私聊类型消息 -- 消息类型处理器
	case common.ChatMsgP2pMsgType:
		msgProcessor := MsgProcessor{
			Conn: this.Conn,
		}
		msgProcessor.processP2pMsg(msg)
	default:
		fmt.Println("未知的消息体类型,无法处理!")
	}
}
