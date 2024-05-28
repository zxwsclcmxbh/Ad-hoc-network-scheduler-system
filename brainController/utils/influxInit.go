package utils

import (
	"log"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

var InfluxClient client.Client
var WriteCh chan *client.Point

func InfluxInit() {
	var err error
	InfluxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: Config.InfluxDB,
	})
	if err != nil {
		log.Fatalln(err)
	}
	WriteCh = make(chan *client.Point)
	go WriteToDB()
}

func WriteToDB() {
	var t []*client.Point
	for p := range WriteCh {
		// log.Println("get from channel", p)
		if len(t) == 10 {
			b, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: "edge", Precision: "ns"})
			// log.Println("batched point", err)
			b.AddPoints(t)
			if err := InfluxClient.Write(b); err != nil {
				log.Println(err)
			}
			t = []*client.Point{p}
		} else {
			t = append(t, p)
		}
	}
}
