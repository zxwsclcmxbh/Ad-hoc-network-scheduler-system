package delay

import (
	"fmt"
	"net"
	"time"
)

const message = "qwertyuiopasdfghjklzxcbvbnm"

func DelayClient(ip, port string) (string, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	start := time.Now()
	conn.Write([]byte(message))
	_, err = conn.Read(buffer)
	if err != nil {
		return "", err
	}
	end := time.Now()
	elapsed := end.Sub(start)
	return fmt.Sprint(elapsed), nil
}
