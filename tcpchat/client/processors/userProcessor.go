package processors

import (
	"awesomeProject/tcpchat/common"
	"encoding/json"
	"fmt"
	"net"
)

// 将方法与结构体进行绑定,以面向对象方式调用
type UserProcessor struct {
	Conn net.Conn
}

// 连接到服务器
func connectToServer() (conn net.Conn, err error) {
	// 连接到服务器
	conn, err = net.Dial("tcp", "localhost:8888")
	return
}

// 协程任务: 和服务器通信
func processConnect(conn net.Conn) {
	// 延时关闭
	defer conn.Close()
	// 创建消息处理器
	processor := &Processor{
		Conn: conn,
	}
	err := processor.ProcessServerConn()
	if err != nil {
		fmt.Errorf("read packet error[%v]!\n", err)
		// 结束协程
		return
	}
}

// 登录到服务器
func (this *UserProcessor) DoLogin() (err error) {
	// 连接到服务器
	conn, err := connectToServer()
	if err != nil {
		fmt.Errorf("Connect to server error[%v]!\n", err)
		return
	}
	defer conn.Close()
	// 登录
	var userId int
	var password string
	fmt.Println("请输入用户ID: ")
	fmt.Scanln(&userId)
	fmt.Println("请输入用户密码: ")
	fmt.Scanln(&password)
	// 构造登录请求体,并序列化,封装成Message,再序列化,再发送
	loginBody := common.LoginBody{
		userId,
		password,
	}
	transfer := &common.Transfer{
		Conn: conn,
	}
	err = transfer.ConvertToMessageAndSend(loginBody, common.LoginMsgType)
	if err != nil {
		fmt.Errorf("Login failed: [%v]!\n", err)
		return
	}
	// 处理服务器的登录响应
	msg, err := transfer.RecvDataPacket()
	if err != nil {
		fmt.Errorf("Login failed: client recv data packet error[%v]!\n", err)
		return
	}
	// 反序列化得到 响应体
	var respBody common.LoginResp
	err = json.Unmarshal([]byte(msg.Data), &respBody)
	if err != nil {
		fmt.Errorf("Login failed: json unmarshal error[%v]!\n", err)
		return
	}
	if respBody.Code != 200 {
		fmt.Println(respBody.Msg)
		return
	} else {
		// 登录成功
		fmt.Printf("------------------恭喜用户[%d]登录成功!!!------------------\n", userId)
		// 保存服务器返回的在线列表
		GlobalConnManager.AddUserBatch(respBody.OnlineList)
		// 保存当前客户的其他信息
		GlobalConnManager.UserId = userId
		GlobalConnManager.Conn = conn
		// 启动协程与服务端通信
		go processConnect(conn)
		// 显示主界面,这里面是个死循环,所以要写开启协程,不然阻塞了,它下面的代码无法运行
		ShowMenu()
	}
	return
}

/**
 * 处理注册
 * 注册请求 最好重新建立一个连接,注册完成后断开
 * 因为登录成功后会和服务器建立连接,一直接收服务器端的消息,那如果,我选择退出登录后,选择注册,
 *			如果没重新连接,那就继续用之前这个,这时服务器的返回消息会被之前那个一直接收服务器消息部分逻辑获取到,那我注册方法里面接收不到回复了
 *          如果重新连接,那么注册就单独开了,我能接收到服务器的回复,注册成功后这个连接就断开了.就和其他部分的逻辑独立开了.不会混乱.
 * 实际上,正常逻辑是先注册,再登录,再长连接通信,
 * 我们只是为了避免.退出登录后没有关闭连接,此时继续去执行注册,就会导致此时的服务器回复被之前的通信部分逻辑获取到,就不好处理了
 */
func (this *UserProcessor) DoRegister() (err error) {

	// 连接到服务器
	conn, err := connectToServer()
	if err != nil {
		fmt.Errorf("Connect to server error[%v]!\n", err)
		return
	}
	defer conn.Close()

	var userId int
	var password string
	var username string
	fmt.Println("请输入用户ID: ")
	fmt.Scanln(&userId)
	fmt.Println("请输入用户密码: ")
	fmt.Scanln(&password)
	fmt.Println("请输入用户名(昵称):")
	fmt.Scanln(&username)
	// 1.构造注册请求体,并序列化
	registerBody := common.RegisterBody{
		userId,
		password,
		username,
	}
	// 3.发送注册请求给服务器
	transfer := &common.Transfer{
		Conn: conn,
	}
	err = transfer.ConvertToMessageAndSend(registerBody, common.RegisterMsgType)
	if err != nil {
		fmt.Errorf("Register failed: [%v]\n", err)
		return
	}
	// 4.处理服务器的注册响应
	msg, err := transfer.RecvDataPacket()
	if err != nil {
		fmt.Errorf("Register failed: client recv data packet error[%v]!\n", err)
		return
	}
	// 反序列化得到 响应体
	var respBody common.RegisterResp
	err = json.Unmarshal([]byte(msg.Data), &respBody)
	if err != nil {
		fmt.Errorf("Register failed: json unmarshal error[%v]!\n", err)
		return
	}
	if respBody.Code != 200 {
		fmt.Println(respBody.Msg)
		return
	} else {
		// 注册成功
		fmt.Printf("用户[%d]注册成功,快快登录进来吧!\n", userId)
	}
	return
}

// 收到用户在线状态变更消息
func (this *UserProcessor) ProcessNotifyOnlineStatusResp(msg *common.Message) (err error) {
	// 反序列化得到 响应体
	var respBody common.UserOnlineStatusBody
	err = json.Unmarshal([]byte(msg.Data), &respBody)
	if err != nil {
		fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	// 某个用户登录了
	if respBody.Status == common.UserStatusOnline {
		GlobalConnManager.AddUser(respBody.UserId)
		// 某个用户退出了
	} else if respBody.Status == common.UserStatusOffline {
		GlobalConnManager.DelUser(respBody.UserId)
	}
	return
}
