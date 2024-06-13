package db

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

const InfluxDBAddr = "http://10.112.123.250:8086"

func WriteInfluxDB(name string, tags map[string]string, fields map[string]interface{}) {
	// 初始化 InfluxDB 客户端
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: InfluxDBAddr,
	})

	if err != nil {
		log.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	pt, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// 写入数据点到 InfluxDB
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "mydb",
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)
	err = c.Write(bp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// 查询 InfluxDB 数据点
	q := client.Query{
		Command:  `SELECT "value" FROM "temperature"`,
		Database: "mydb",
	}
	res, err := c.Query(q)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	var value float64
	for _, row := range res.Results[0].Series[0].Values {
		value = row[1].(float64)
	}
	fmt.Printf("%+v", value)
}
