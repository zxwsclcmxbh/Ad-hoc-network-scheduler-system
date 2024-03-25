package delay

import (
	"fmt"
	"net"
	"time"
)

func DelayClient(ip, port string) string {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer conn.Close()

	message := "qwertyuiopasdfghjklzxcbvbnm"
	buffer := make([]byte, 1024)
	start := time.Now()
	conn.Write([]byte(message))
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	end := time.Now()
	elapsed := end.Sub(start)
	return fmt.Sprint(elapsed)
}
