package main

import (
	"awesomeProject/tcpchat/common"
	"awesomeProject/tcpchat/server/dao"
	"awesomeProject/tcpchat/server/processors"
	"fmt"
	"io"
	"net"
)

// 协程任务: 监听到客户端连接后，处理
func processConnect(conn net.Conn) {
	defer conn.Close()
	// 创建处理器,处理与客户端的连接
	processor := &processors.Processor{
		Conn: conn,
	}
	err := processor.ProcessClientConn()
	if err != nil {
		if err == io.EOF {
			fmt.Println("客户端退出,服务器协程正常结束!")
			// 通知其他用户此人下线
			// 从管理列表中移除此连接
			userProcessor := processors.GlobalConnManager.GetProcessorByConn(conn)
			if userProcessor != nil {
				userProcessor.NotifyOtherMyOnlineStatus(common.UserStatusOffline)
				processors.GlobalConnManager.DelUserByUid(userProcessor.UserId)
			}
			return
		} else {
			fmt.Errorf("read packet error[%v]!\n", err)
		}
		// 结束协程
		return
	}
}

func RunServer() {
	fmt.Println("********服务器开始在8888端口监听*******")
	// 初始化redis
	dao.InitRedisClient()
	// 监听
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Errorf("server listen error[%v]!\n", err)
	}
	defer listener.Close()
	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("server accept error[%v]!\n", err)
		}
		fmt.Printf(">>>>>>>>一个客户端已连接[%v]>>>>>>>>>\n", conn.RemoteAddr().String())
		// 启动协程,处理连接
		go processConnect(conn)
	}

}

func main() {
	RunServer()
}
