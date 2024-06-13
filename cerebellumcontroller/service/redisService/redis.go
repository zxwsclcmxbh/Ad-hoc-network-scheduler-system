package redisService

import (
	"context"
	"encoding/json"
	"fmt"
	"iai/cerebellumController/config"
	"iai/cerebellumController/dto/message"
	"github.com/go-redis/redis/v8"
)

var Conn *redis.Client

func Init() {
	Conn = redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
		DB:   0})
}
func Publish(prefix string, name string, msg *message.Message) error {
	ctx :=context.Background()
	data, _ := json.Marshal(msg)
	err := Conn.Publish(ctx, fmt.Sprintf("/%s/%s", prefix, name), string(data)).Err()
	return err
}

func GetSubHandler(prefix string, name string) *redis.PubSub {
	ctx := context.Background()
	return Conn.Subscribe(ctx, fmt.Sprintf("/%s/%s", prefix, name))
}
