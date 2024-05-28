package worker

import (
	"context"
	"encoding/json"
	"iai/cerebellumController/dto/message"
	"iai/cerebellumController/service/outgoingMessage"
	"iai/cerebellumController/service/redisService"
	"iai/cerebellumController/util"
	"log"
)

func HandleOutgoingMessage(prefix string, name string, node string, dest string, svc string, ctx context.Context) {
	log.Println(prefix, name, node, dest, svc)
	sub := redisService.GetSubHandler(prefix, name)
	ch := sub.Channel()
	defer sub.Close()
	for {
		select {
		case msg_redis := <-ch:
			msg_str := msg_redis.Payload
			msg := message.Message{}
			json.Unmarshal([]byte(msg_str), &msg)
			if err := outgoingMessage.SendMessage(util.GenHost(svc, node), &msg, name, dest, prefix); err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
