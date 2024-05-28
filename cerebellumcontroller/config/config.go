package config

import (
	"log"

	"github.com/jinzhu/configor"
)

var Config = struct {
	DB struct {
		File string `default:"route.sqlite" env:"dbfile"`
	}
	APP struct {
		Addr string `default:"0.0.0.0" env:"addr"`
		Port string `default:"3000" env:"port"`
	}
	Storage struct {
		BasePath string `default:"./data/" env:"filebase"`
	}
	Redis struct {
		Addr string `default:"redis.svc:6379" env:"redisaddr"`
	}
	Host  string `default:"node0" env:"nodename"`
	Brain string `default:"ws://192.168.1.3:8080/api/brainController/ws" env:"brainaddr"`
	Debug string `default:"false" env:"debug"`
}{}

func Init() {
	configor.Load(&Config)
	log.Println(Config)
}
