package utils

import (
	"log"

	"github.com/jinzhu/configor"
)

type KV struct {
	Key   string
	Value string
}

var Config = struct {
	DB       string
	NodeIP   []KV
	Registry struct {
		IP   string
		Port string
	}
	GPUNode     string
	ImageMap    []KV
	PHMEndpoint string
	InfluxDB    string
	Minio       struct {
		URL      string
		Username string
		Password string
	}
	Ingress struct {
		BaseUrl  string
		AdminUrl string
		ApiKey   string
	}
}{}

func LoadConf() {
	configor.Load(&Config, "config.yaml")
	log.Println("config", Config)
}
