package main

import (
	"awesomeProject/tcpchat/client/processors"
	"fmt"
)

func main() {
	loop := true
	key := -1
	for loop {
		fmt.Println("------------------欢迎登录多人聊天系统--------------------")
		fmt.Println("\t\t\t1.登录聊天室")
		fmt.Println("\t\t\t2.注册新用户")
		fmt.Println("\t\t\t3.退出系统")
		fmt.Println("请选择(1-3):")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("你选择登录聊天室")
			// 创建处理器对象
			userProcessor := &processors.UserProcessor{}
			userProcessor.DoLogin()
		case 2:
			fmt.Println("你选择注册新用户")
			// 创建处理器对象
			userProcessor := &processors.UserProcessor{}
			userProcessor.DoRegister()
		case 3:
			fmt.Println("你选择退出系统")
			loop = false
		default:
			fmt.Println("你的输入不正确,请重试!")
		}
	}
}
