package processors

import (
	"awesomeProject/tcpchat/common"
	"awesomeProject/tcpchat/server/dao"
	"awesomeProject/tcpchat/server/errors"
	"awesomeProject/tcpchat/server/model"
	"encoding/json"
	"fmt"
	"net"
)

// 将方法与结构体进行绑定,以面向对象方式调用
type UserProcessor struct {
	// 这是和哪个用户的连接
	Conn net.Conn
	// 这个用户的ID
	UserId int
}

// 处理登录请求
func (this *UserProcessor) ProcessLoginReq(msg *common.Message) {
	// 反序列化得到消息体a
	var loginBody common.LoginBody
	err := json.Unmarshal([]byte(msg.Data), &loginBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 判断登录信息并返回登录结果
	// 1.先构造返回消息体的实际内容，并序列化
	var respBody common.LoginResp
	userDAO := &dao.UserDAO{}
	// 判断用户是否已经登录
	onlineList := GlobalConnManager.GetAllUser()
	flag := false
	for _, id := range onlineList {
		if id == loginBody.UserId {
			flag = true
			break
		}
	}
	if flag {
		respBody.Code = 400
		respBody.Msg = "您已登录!,请勿重复操作!"
	} else {
		// 调用持久层方法判断用户是否存在
		_, err = userDAO.Login(loginBody.UserId, loginBody.Password)
		if err == errors.USER_NOT_EXIST || err == errors.INVALID_UID_PASSWD {
			respBody.Code = 400
			// 取出错误信息
			respBody.Msg = err.Error()
		} else if err != nil {
			respBody.Code = 500
			respBody.Msg = "失败了,服务器冒烟了..."
		} else {
			respBody.Code = 200
			respBody.Msg = "登录成功!"
			// 告诉他在线列表
			respBody.OnlineList = GlobalConnManager.GetAllUser()
		}
	}
	// 序列化消息体->封装成Message->序列化并发送
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	transfer.ConvertToMessageAndSend(respBody, common.LoginRespMsgType)
	// ---------登录成功后-----------
	if respBody.Code == 200 {
		// 如果是登录成功，那么把他加入管理列表，
		this.UserId = loginBody.UserId
		GlobalConnManager.AddNewUser(loginBody.UserId, this)
		// 还要通知所有人我上线了
		this.NotifyOtherMyOnlineStatus(common.UserStatusOnline)
	}
}

// 通知其他人我上线了
func (this *UserProcessor) NotifyOtherMyOnlineStatus(myOnlineStatus int) {
	// 先拿到当前全部在线列表
	userSlice := GlobalConnManager.GetAllUser()
	for _, userId := range userSlice {
		// 避开自己
		if userId == this.UserId {
			continue
		}
		// 拿到这个人的处理器
		up := GlobalConnManager.GetProcessorByUserId(userId)
		// 告诉这个人我上线了
		up.ProcessUpdateOnlineListReq(this.UserId, myOnlineStatus)
	}
}

// 处理请求:告诉这个客户端,某个人的在线状态改变
func (this *UserProcessor) ProcessUpdateOnlineListReq(userId int, newStatus int) {
	// 构造响应体
	respBody := common.UserOnlineStatusBody{
		UserId: userId,
		Status: newStatus,
	}
	// 序列化消息体->封装成Message->序列化并发送
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	transfer.ConvertToMessageAndSend(respBody, common.NotifyUserStatusMsgType)
}

// 处理注册请求
func (this *UserProcessor) ProcessRegisterReq(msg *common.Message) {
	// 反序列化得到消息体
	var registerBody common.RegisterBody
	err := json.Unmarshal([]byte(msg.Data), &registerBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 判断登录信息并返回登录结果
	// 1.先构造返回消息体的
	var respBody common.RegisterResp
	// 这一步相当于把前端数据转为后端实体
	user := &model.User{registerBody.UserId, registerBody.Password, registerBody.Username}
	// 调用持久层方法
	userDAO := &dao.UserDAO{}
	err = userDAO.Register(user)
	if err == errors.USER_ALREADY_EXIST {
		respBody.Code = 400
		respBody.Msg = "用户ID已存在!"
	} else if err != nil {
		respBody.Code = 500
		respBody.Msg = "注册失败！"
	} else {
		respBody.Code = 200
		respBody.Msg = "注册成功！"
	}
	// 序列化消息体->封装成Message->序列化并发送
	transfer := &common.Transfer{
		Conn: this.Conn,
	}
	transfer.ConvertToMessageAndSend(respBody, common.RegisterRespMsgType)
}
