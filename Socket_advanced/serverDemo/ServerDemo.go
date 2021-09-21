package main

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	layout = "2006-01-02 15:04:05"
	Kitchen     = "3:04PM"
	prefixName = "client"
)
var (
	suffixName int32 = 1
	// 消息发送给在线用户
	onlineUser = make(map[string]client)
	msgList = make(chan []byte)
)

type client struct {
	msgChan chan string
	name string
	addr net.Addr
}

// 发送消息给用户
func sendMsgToUser(user client, con net.Conn) {
	for {
		for v := range user.msgChan {
			con.Write([]byte(v))
		}
	}
}

// 监听msgList通道中的消息将消息发给所有client
func listenMsg() {
	for {
		msg := <- msgList
		for _, client := range onlineUser {
			client.msgChan <- string(msg)
		}
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	//userStatus := make(chan bool)
	quitList := make(chan bool)
	// 创建user将齐放入在线表中
	user := client{
		msgChan: make(chan string),
		name: prefixName + strconv.Itoa(int(suffixName)),
		addr: conn.RemoteAddr(),
	}
	fmt.Println(user.name, "已连上服务器")
	// suffixName ++
	atomic.AddInt32(&suffixName, 1)
	// ip地址是唯一的不会出现并发写情况
	onlineUser[conn.RemoteAddr().String()] = user

	go sendMsgToUser(user, conn)
	go func() {
		var tmp [128]byte
		for {
			read, err := conn.Read(tmp[:])
			if read == 0 {
				// 通知要客户端打印到客户端
				quitList<-true
				return
			}
			if err != nil {
				fmt.Println("read from conn failed, err:", err)
				return
			}
			msgList <- []byte(user.name + " " + time.Now().Format(Kitchen) + ":" + string(tmp[:read]))
			fmt.Println(time.Now().Format(layout), ":", string(tmp[:read]))
			//userStatus<-true
		}
	}()

	for  {
		select {
		// 进行相关操作，并退出该连接
		case <-quitList:
			// 将用户退出消息打印在终端
			fmt.Println(conn.RemoteAddr().String(),time.Now().Format(layout), "已下线!")
			// 将下线用户从online表中删除
			delete(onlineUser, conn.RemoteAddr().String())
			return
		//case <-userStatus:
		case <-time.After(time.Second * 20):
			// 超过20秒断开则不是活跃用户，不再接受信息
			delete(onlineUser, user.addr.String())
			fmt.Println(user.name, "超时退出!")
			msgList <- []byte(fmt.Sprintf("%s 超时退出!", user.name))
			return
		}
	}
}
func main() {
	listen, err := net.Listen("tcp", "localhost:8089")
	if err != nil {
		fmt.Println("start server on failed ", err)
	}
	go listenMsg()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err: ", err)
			return
		}
		go process(conn)
	}
}
