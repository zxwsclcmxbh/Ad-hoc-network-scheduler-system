package router

import (
	messageHandler "iai/cerebellumController/controller/messageHandler"

	"github.com/gofiber/fiber/v2"
)

func InitMessageROuter(app *fiber.App) {
	msg := app.Group("/api").Group("/v1").Group("/message")
	msg.Post("/incoming", messageHandler.HandleIncoming)

}
