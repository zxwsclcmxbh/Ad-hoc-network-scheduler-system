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

var Source = map[string]string{
	"edge-node-1": "192.168.1.1",
	"edge-node-2": "192.168.1.2",
	"edge-node-3": "192.168.1.3",
	"edge-node-4": "192.168.1.4",
}

var endpoint = "http://10.112.123.250:8080/reportNetwork"

const timeInterval = 1
const containerID = "frr"
const Port = "8100"

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

func main() {
	name, _ := os.Hostname()
	res := &NetworkData{
		Name: name,
		IP:   Source[name],
	}

	for {
		time.Sleep(timeInterval * time.Second)
		nsms, err := nsm.GetNsm(containerID)
		if err != nil {
			break
		}
		log.Println("nsm:", nsms)
		res.Nsms = []Nsm{}
		for i := 0; i < len(nsms); i++ {
			delayRes, err := delay.DelayClient(nsms[i], Port)
			if err != nil {
				break
			}
			nsm := &Nsm{
				Name:  getKeyByValue(Source, nsms[i]),
				IP:    nsms[i],
				Delay: delayRes,
			}
			res.Nsms = append(res.Nsms, *nsm)
		}
		err = sendPostRequest(endpoint, res)
		if err != nil {
			log.Println("err:", err)
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
