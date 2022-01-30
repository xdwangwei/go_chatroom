package processors

import (
	"awesomeProject/tcpchat/common"
	"awesomeProject/tcpchat/server/dao"
	"awesomeProject/tcpchat/server/errors"
	"encoding/json"
	"fmt"
	"net"
)

type MsgProcessor struct {
	Conn net.Conn
}

// 处理接收到群消息
func (this *MsgProcessor) processGroupMsg(msg *common.Message) (err error) {
	// 反序列化得到 响应体
	var msgBody common.ChatMsgBody
	err = json.Unmarshal([]byte(msg.Data), &msgBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 可以通过 消息体里面的 senderId 判断发送者,但为了避免客户端传假数据,这里根据 此条连接 来判断真实来源
	up := GlobalConnManager.GetProcessorByConn(this.Conn)
	// 发送者
	senderId := up.UserId
	// 将消息发送给所有在线用户
	userList := GlobalConnManager.GetAllUser()
	// 遍历在线列表
	for _, uid := range userList {
		// 跳过发送者
		if uid == senderId {
			continue
		}
		// 拿到这个人的处理器
		up := GlobalConnManager.GetProcessorByUserId(uid)
		// 构造消息处理器
		// 将消息发过去
		mp := MsgProcessor{
			Conn: up.Conn,
		}
		mp.forwardMsg(senderId, msgBody.Content, common.ChatMsgGroupMsgType)
	}
	return
}

// 处理接收到私聊消息
func (this *MsgProcessor) processP2pMsg(msg *common.Message) (err error) {
	// 反序列化得到 响应体
	var msgBody common.ChatMsgBody
	err = json.Unmarshal([]byte(msg.Data), &msgBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 可以通过 消息体里面的 senderId 判断发送者,但为了避免客户端传假数据,这里根据 此条连接 来判断真实来源
	up := GlobalConnManager.GetProcessorByConn(this.Conn)
	// 发送者
	senderId := up.UserId
	// 接收者
	receiverId := msgBody.ReceiverId

	// 构造响应体,告诉客户端此条消息的处理结果
	respBody := common.ChatMsgResp{}
	respBody.Code = 400
	// 将原消息体设置进去
	respBody.Data = msg.Data

	// 不允许自己给自己发送消息
	if senderId == receiverId {
		respBody.Msg = "不能自己给自己发送消息!"
	} else {
		// 判断接收者是否存在
		_, err = dao.GlobalUserDAO.GetUserById(receiverId)
		// 接收者不存在
		if err == errors.USER_NOT_EXIST {
			respBody.Msg = "此用户不存在!"
		} else {
			// 判断接收者是否在线
			online := false
			userList := GlobalConnManager.GetAllUser()
			for _, id := range userList {
				if id == receiverId {
					online = true
					break
				}
			}
			if !online {
				// 用户不在线
				respBody.Msg = "此用户并不在线,暂不支持离线留言!"
			} else {
				// 用户在线
				respBody.Code = 200
				respBody.Msg = "已转发给目标用户"
				// 1.转发到目标用户
				// 拿到这个人的处理器
				up := GlobalConnManager.GetProcessorByUserId(receiverId)
				// 构造消息处理器
				// 将消息发过去
				mp := MsgProcessor{
					Conn: up.Conn,
				}
				mp.forwardMsg(senderId, msgBody.Content, common.ChatMsgP2pMsgType)
			}
		}
	}
	// 告诉发送者结果
	// 序列化消息体->封装成Message->序列化并发送
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	transfer.ConvertToMessageAndSend(respBody, common.ChatMsgRespType)
	return
}

// 转发消息消息给客户
func (this *MsgProcessor) forwardMsg(senderId int, content string, msgType string) {
	// 构造响应体
	respBody := common.ChatMsgBody{
		SenderId: senderId,
		Content:  content,
	}
	// 序列化消息体->封装成Message->序列化并发送
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	transfer.ConvertToMessageAndSend(respBody, msgType)
}
