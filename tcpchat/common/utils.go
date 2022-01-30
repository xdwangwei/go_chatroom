package common

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

// 将方法与结构体进行绑定,以面向对象方式调用
type Transfer struct {
	Conn net.Conn
	Buf  [4096]byte
}

// 发送 消息
// 每个消息都是 先发送此消息的长度(用四个字节)，再发送此消息的内容
func (this *Transfer) SendDataPacket(bytes []byte) (err error) {
	// 先发送消息长度,用4字节
	msgLen := len(bytes)
	binary.LittleEndian.PutUint32(this.Buf[:4], uint32(msgLen))
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		//fmt.Errorf("write data packet error!\n")
		return
	}
	// 再发送消息体
	n, err = this.Conn.Write(bytes)
	if n != msgLen || err != nil {
		//fmt.Errorf("write data packet error!\n")
		return
	}
	return
}

// 接收并解析消息
func (this *Transfer) RecvDataPacket() (msg Message, err error) {
	n, err := this.Conn.Read(this.Buf[:4])
	// 第一次接收到的是消息的长度，应该是四个字节
	if n != 4 || err != nil {
		//fmt.Errorf("read data packet error!\n")
		return
	}
	// 转换为整形
	msgLen := binary.LittleEndian.Uint32(this.Buf[:4])
	// 再接受消息体
	n, err = this.Conn.Read(this.Buf[:])
	// 实际接收到的消息体的大小和刚才接收到的客户端给的大小不一致
	if n != int(msgLen) || err != nil {
		//fmt.Errorf("read data packet error!\n")
		return
	}
	// 消息体反序列化
	// 使用命名返回值,不用创建变量
	err = json.Unmarshal(this.Buf[:msgLen], &msg)
	if err != nil {
		//fmt.Errorf("json unmarshal error[%v]!\n", err)
		return
	}
	return
}

// 把一个结构体序列化后封装成一个Message结构体,再序列化,再发送出去
// messageType是标记封装成哪一类messag
// 没有进行类型校验
func (this *Transfer) ConvertToMessageAndSend(st interface{}, msgType string) (err error) {
	// 1.序列化消息内容
	respBytes, err := json.Marshal(&st)
	if err != nil {
		//fmt.Errorf("json marshal error[%v]!\n", err)
		return
	}
	// 2.将其封装成 Message 类型
	respMsg := Message{
		msgType,
		string(respBytes),
	}
	// 3.将 message 对象序列化 并发送给客户端
	msgBytes, err := json.Marshal(&respMsg)
	if err != nil {
		//fmt.Errorf("json marshal error[%v]!\n", err)
		return
	}
	// 4.发送
	err = this.SendDataPacket(msgBytes)
	if err != nil {
		//fmt.Errorf("server send data packet error[%v]!\n", err)
		return
	}
	return
}
