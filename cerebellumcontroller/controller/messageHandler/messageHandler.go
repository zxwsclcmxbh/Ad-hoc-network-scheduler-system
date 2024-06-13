package messagehandler

import (
	"iai/cerebellumController/config"
	"iai/cerebellumController/dto/common"
	"iai/cerebellumController/dto/message"
	"iai/cerebellumController/service/redisService"
	"iai/cerebellumController/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

func HandleIncoming(c *fiber.Ctx) error {
	incoming := &message.MessageIncoming{}
	if err := c.BodyParser(incoming); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	}
	log.Println(incoming)
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	}
	fname := util.GenFileName(config.Config.Storage.BasePath, file.Filename, incoming.TaskId)
	c.SaveFile(file, fname)
	msg := message.Message{File: util.ReGenFileNameIncoming(fname, incoming.TaskId), Length: 0, MetaData: incoming.MetaData}
	redisService.Publish(incoming.TaskId, incoming.Dst, &msg)
	return c.JSON(common.Response{Code: 0, Msg: "ok"})
}
