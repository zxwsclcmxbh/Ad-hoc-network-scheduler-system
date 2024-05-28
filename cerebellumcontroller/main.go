package main

import (
	"iai/cerebellumController/config"
	"iai/cerebellumController/dao"
	"iai/cerebellumController/router"
	"iai/cerebellumController/service/brain"
	"iai/cerebellumController/service/redisService"
	"iai/cerebellumController/service/routerService"
	"iai/cerebellumController/wsServer"
	"os"

	"github.com/gofiber/fiber/v2"
)

func init() {
	config.Init()
	dao.Init()
	redisService.Init()
	routerService.Init()
	go wsServer.Run(make(chan os.Signal, 1))
	go brain.SendWorker()
}

func main() {
	app := fiber.New()
	router.InitControlRouter(app)
	router.InitMessageROuter(app)
	router.InitReport(app)
	app.Listen(config.Config.APP.Addr + ":" + config.Config.APP.Port)
}
