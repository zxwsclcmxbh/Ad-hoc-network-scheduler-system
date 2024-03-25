package main

import (
	"fmt"
	"net"
)

func DelayServer(ip, port string) {
	l, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		message := string(buffer[:n])
		conn.Write([]byte(message))
	}
}
func main() {
	DelayServer("0.0.0.0", "8100")
}
