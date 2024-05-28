package main

import (
	"cloud/db"
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
	// router := gin.Default()

	// router.POST("/reportNetwork", handlePost)

	// fmt.Println("Server is running on port 8080...")
	// router.Run(":8080")
	networkData := NetworkData{
		Name: "edge-node-1",
		IP:   "192.168.1.1",
		Nsms: []Nsm{
			{
				Name:    "edge-node-2",
				IP:      "192.168.1.2",
				Delay:   "3.1476s",
				CalTime: "1716019020861",
			},
			{
				Name:    "edge-node-3",
				IP:      "192.168.1.3",
				Delay:   "0.1476s",
				CalTime: "1716019020863",
			},
			{
				Name:    "edge-node-4",
				IP:      "192.168.1.4",
				Delay:   "4.1476s",
				CalTime: "1716019020867",
			},
		},
	}
	fmt.Println("Received NetworkData:")
	fmt.Println("Name:", networkData.Name)
	fmt.Println("IP:", networkData.IP)
	fmt.Println("Nsms:")
	for _, nsm := range networkData.Nsms {
		fmt.Println("  Name:", nsm.Name)
		fmt.Println("  IP:", nsm.IP)
		fmt.Println("  Delay:", nsm.Delay)
	}
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
		// 写入数据到InfluxDB
		tags := map[string]string{
			"src": networkData.Name,
			"dst": nsm.Name,
		}
		fields := map[string]interface{}{
			"value": nsm.Delay,
		}
		db.WriteInfluxDB("delay", tags, fields)
	}
	c.JSON(200, gin.H{"message": "Received NetworkData"})
}
