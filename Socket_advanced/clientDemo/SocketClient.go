package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func writeToServer(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		conn.Write(line)
	}
}

func getMsg(conn net.Conn) {
	var tmp [128]byte
	for {
		n, err := conn.Read(tmp[:])
		if err != nil {
			fmt.Println(conn.LocalAddr(), ":offline")
			return
		}

		fmt.Printf("%s\n", string(tmp[:n]))
		//fmt.Println(conn.RemoteAddr(),":",string(tmp[:n]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8089")
	if err != nil {
		fmt.Println("dial localhost:8087 fail!")
	}
	defer conn.Close()
	go writeToServer(conn)
	getMsg(conn)
}
