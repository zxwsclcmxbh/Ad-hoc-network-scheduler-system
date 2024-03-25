package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type NetworkData struct {
	Name string `json:"name"` // 当前节点名称
	IP   string `json:"ip"`   // 当前节点IP
	Nsms []Nsm  `json:"nsms"` // 邻居信息
}

type Nsm struct {
	Name    string `json:"name"`     // 邻居节点名称
	IP      string `json:"ip"`       // 邻居节点IP
	Delay   string `json:"delay"`    // 时延
	CalTime string `json:"cal_time"` // 计算时间
}

func main() {
	router := gin.Default()

	router.POST("/reportNetwork", handlePost)

	fmt.Println("Server is running on port 8080...")
	router.Run(":8080")
}

func handlePost(c *gin.Context) {
	var networkData NetworkData
	if err := c.BindJSON(&networkData); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	// 打印接收到的数据
	fmt.Println("Received NetworkData:")
	fmt.Println("Name:", networkData.Name)
	fmt.Println("IP:", networkData.IP)
	fmt.Println("Nsms:")
	for _, nsm := range networkData.Nsms {
		fmt.Println("  Name:", nsm.Name)
		fmt.Println("  IP:", nsm.IP)
		fmt.Println("  Delay:", nsm.Delay)
	}

	c.JSON(200, gin.H{"message": "Received NetworkData"})
}
