package processors

import (
	"awesomeProject/tcpchat/common"
	"encoding/json"
	"fmt"
)

// 负责处理 聊天类型的消息
type MsgProcessor struct {
	// Conn 可以通过全局对象ConnManager获取到,所以这里不用传了
	// Conn net.Conn
}

// 发送群聊类型消息
func (this *MsgProcessor) SendGroupMsg(content string) (err error) {
	// 构造请求体,并序列化,封装成Message,再序列化,再发送
	msgBody := common.ChatMsgBody{
		SenderId: GlobalConnManager.UserId,
		Content:  content,
	}
	transfer := &common.Transfer{
		Conn: GlobalConnManager.Conn,
	}
	err = transfer.ConvertToMessageAndSend(msgBody, common.ChatMsgGroupMsgType)
	if err != nil {
		fmt.Errorf("send group message error[%v]!\n", err)
		return
	}
	return
}

// 发送私聊类型消息
func (this *MsgProcessor) SendP2pMsg(receiverId int, content string) (err error) {
	// 构造请求体,并序列化,封装成Message,再序列化,再发送
	msgBody := common.ChatMsgBody{
		SenderId:   GlobalConnManager.UserId,
		ReceiverId: receiverId,
		Content:    content,
	}
	transfer := &common.Transfer{
		Conn: GlobalConnManager.Conn,
	}
	err = transfer.ConvertToMessageAndSend(msgBody, common.ChatMsgP2pMsgType)
	if err != nil {
		fmt.Errorf("send p2p message error[%v]!\n", err)
		return
	}
	return
}

// 处理 服务端转发过来的群聊/私聊消息
func (this *MsgProcessor) ProcessMsgFromServer(msg *common.Message) (err error) {
	// 反序列化得到 响应体
	var respBody common.ChatMsgBody
	err = json.Unmarshal([]byte(msg.Data), &respBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 如果是群聊消息
	if msg.MessageType == common.ChatMsgGroupMsgType {
		fmt.Printf("%d 对大家说: %s\n", respBody.SenderId, respBody.Content)
		// 如果是私聊消息
	} else if msg.MessageType == common.ChatMsgP2pMsgType {
		fmt.Printf("%d 私聊我: %s\n", respBody.SenderId, respBody.Content)
	}
	return
}

// 处理 服务端的响应,
func (this *MsgProcessor) ProcessMsgRespFromServer(msg *common.Message) (err error) {
	// 反序列化得到 响应体
	var respBody common.ChatMsgResp
	err = json.Unmarshal([]byte(msg.Data), &respBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 若消息发送失败
	if respBody.Code != 200 {
		// 尝试反序列化
		var msgBody common.ChatMsgBody
		err := json.Unmarshal([]byte(respBody.Data), &msgBody)
		// 反序列化失败就以字符串形式输出
		if err != nil {
			fmt.Printf("消息发送失败[%s],原消息数据包: %s\n", respBody.Msg, respBody.Data)
		} else {
			fmt.Printf("消息发送失败[%v],原消息数据包: %v\n", respBody.Msg, msgBody)
		}
	}
	return
}
