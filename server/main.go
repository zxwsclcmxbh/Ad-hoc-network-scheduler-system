package main

import (
	"net"
	"os"

	"log"
)

var Source = map[string]string{
	"edge-node-1": "192.168.1.1",
	"edge-node-2": "192.168.1.2",
	"edge-node-3": "192.168.1.3",
	"edge-node-4": "192.168.1.4",
}

const Port = "8100"

func DelayServer(ip string) {
	l, err := net.Listen("tcp", ip+":"+Port)
	if err != nil {
		log.Println("err:", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("err:", err)
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
	name, _ := os.Hostname()
	log.Println("delay server:", Source[name]+":"+Port)
	DelayServer(Source[name])
}
