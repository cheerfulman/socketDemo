package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(line)
}

func getMsg(conn net.Conn) {
	var tmp [128]byte
	for {
		n, err := conn.Read(tmp[:])
		if err != nil {
			fmt.Println(conn.LocalAddr(), ":client read fail")
		}

		fmt.Printf("%s\n", string(tmp[:n]))
		//fmt.Println(conn.RemoteAddr(),":",string(tmp[:n]))
	}
}
func main() {
	conn, err := net.Dial("tcp", "localhost:8087")
	if err != nil {
		fmt.Println("dial localhost:8087 fail!")
	}
	defer conn.Close()
	for {
		conn.Write([]byte(getInput()))
		go getMsg(conn)
	}
}
