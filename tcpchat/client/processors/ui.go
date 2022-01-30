package processors

import (
	"fmt"
)

// 打印在线用户列表
func printOnlineList() {
	userList := GlobalConnManager.GetAllUser()
	fmt.Println("*********************************************")
	if len(userList) == 0 {
		fmt.Println("\t当前只有您一人在线哦,请耐心等待!\n")
	} else {
		fmt.Println("\t当前在线用户列表如下:")
		for _, id := range userList {
			fmt.Println("\t用户ID: ", id)
		}
	}
	fmt.Println("*********************************************")
}

// 发送群聊消息
func sendGroupMsg() {
	content := ""
	fmt.Println("请输入消息内容:")
	fmt.Scanln(&content)
	// 调用处理器
	mp := &MsgProcessor{}
	mp.SendGroupMsg(content)
}

// 发送私聊消息
func sendP2pMsg() {
	receiverId := -1
	fmt.Println("请输入接收方ID:")
	fmt.Scanln(&receiverId)
	content := ""
	fmt.Println("请输入消息内容:")
	fmt.Scanln(&content)
	// 调用处理器
	mp := &MsgProcessor{}
	mp.SendP2pMsg(receiverId, content)
}

func ShowMenu() {
	loop := true
	key := -1
	for loop {
		fmt.Println("\t\t\t1.显示在线用户列表")
		fmt.Println("\t\t\t2.发送群聊消息")
		fmt.Println("\t\t\t3.发送私聊消息")
		fmt.Println("\t\t\t4.退出登录")
		fmt.Println("请选择(1-4):")
		fmt.Scanln(&key)
		switch key {
		case 1:
			printOnlineList()
		case 2:
			sendGroupMsg()
		case 3:
			sendP2pMsg()
		case 4:
			fmt.Println("已退出登录!")
			loop = false
		default:
			fmt.Println("你的选择有误,请重试!")
		}
	}
}
