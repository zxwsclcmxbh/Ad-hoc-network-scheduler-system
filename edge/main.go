package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"edge/delay"
	"edge/nsm"
)

type NetworkData struct {
	Name string
	IP   string
	Nsms []Nsm
}

type Nsm struct {
	Name  string
	IP    string
	Delay string
}

var Source = map[string]string{
	"edge-node-1": "192.168.1.1",
	"edge-node-2": "192.168.1.2",
	"edge-node-3": "192.168.1.3",
	"edge-node-4": "192.168.1.4",
}

var endpoint = "http://10.112.123.250:8080/reportNetwork"

func main() {
	time.Sleep(1 * time.Second)
	name, _ := os.Hostname()
	res := &NetworkData{
		Name: name,
		IP:   Source[name],
	}

	for {
		time.Sleep(1 * time.Second)
		nsms := nsm.GetNsm("frr")
		log.Printf("nsm:%v ", nsms)
		res.Nsms = []Nsm{}
		for i := 0; i < len(nsms); i++ {
			nsm := &Nsm{
				Name:  getKeyByValue(Source, nsms[i]),
				IP:    nsms[i],
				Delay: delay.DelayClient(nsms[i], "8100"),
			}
			res.Nsms = append(res.Nsms, *nsm)
		}
		err := sendPostRequest(endpoint, res)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		log.Printf("report %v succ.", res)
	}
}

func getKeyByValue(m map[string]string, value string) string {
	for key, val := range m {
		if val == value {
			return key
		}
	}
	return ""
}

func sendPostRequest(url string, data *NetworkData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
