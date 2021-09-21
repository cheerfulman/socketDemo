package main

import (
	"fmt"
	"net"
	"time"
)

const (
	layout = "2006-01-02 15:04:05"
	Kitchen     = "3:04PM"
)
var (
	connList []net.Conn
)
func process(conn net.Conn) {
	var tmp [128]byte
	for {
		read, err := conn.Read(tmp[:])
		if err != nil {
			fmt.Println("read from conn failed, err:", err)
			return
		}
		fmt.Println(time.Now().Format(layout), ":", string(tmp[:read]))
		for _, clientConn := range connList {
			clientConn.Write([]byte(conn.RemoteAddr().String() + " " + time.Now().Format(Kitchen) + ":" + string(tmp[:read])))
		}
	}
}
func main() {
	listen, err := net.Listen("tcp", "localhost:8087")
	if err != nil {
		fmt.Println("start server on failed ", err)
	}
	for {
		conn, err := listen.Accept()
		connList = append(connList, conn)
		fmt.Println(conn.RemoteAddr(), "已连上服务器")
		if err != nil {
			fmt.Println("accept failed, err: ", err)
			return
		}
		go process(conn)
	}
}
