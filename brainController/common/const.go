package common

import (
	"cloud/brainController/utils"
	"log"
)

type BaseResponse struct {
	StatusCode int32       `json:"code"`
	StatusMsg  string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var NodeIP map[string]string

// {
// 	"cloud-node":     "10.112.245.161",
// 	"cloud-node-gpu": "10.112.245.161",
// }

// const HubAddress = "192.168.1.3"
// const HubPort = "5000"

// const GPUNode = "cloud-node-gpu"

var HubAddress string
var HubPort string
var GPUNode string

var CommonImgae map[string]string

// {
// 	"intake-kafka":               "192.168.1.3:5000/iai-phm/common/intake:v4",
// 	"intake-http-form":           "192.168.1.3:5000/iai-phm/common/intake:v4",
// 	"intake-http-json":           "192.168.1.3:5000/iai-phm/common/intake:v4",
// 	"process-rollingwindow":      "192.168.1.3:5000/iai-phm/common/rollingwindow:v5",
// 	"process-transpose":          "192.168.1.3:5000/iai-phm/common/transpose:v4",
// 	"process-feature-extraction": "192.168.1.3:5000/iai-phm/bfd/feature-extraction:v5",
// 	"output":                     "192.168.1.3:5000/iai-phm/common/output:v4",
// }

func TransformConfig() {
	HubAddress = utils.Config.Registry.IP
	HubPort = utils.Config.Registry.Port
	GPUNode = utils.Config.GPUNode
	NodeIP = make(map[string]string)
	for _, item := range utils.Config.NodeIP {
		NodeIP[item.Key] = item.Value
	}
	log.Println("nodeip", NodeIP)
	CommonImgae = make(map[string]string)
	for _, item := range utils.Config.ImageMap {
		CommonImgae[item.Key] = item.Value
	}
	log.Println("imagemap", CommonImgae)
}
