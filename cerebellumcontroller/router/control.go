package router

import (
	"iai/cerebellumController/controller/routerHandler"

	"github.com/gofiber/fiber/v2"
)

func InitControlRouter(app *fiber.App) {
	control := app.Group("/api").Group("/v1").Group("/contorl")
	control.Post("/new", routerHandler.NewRoute)
	control.Post("/modify", routerHandler.ModifyRoute)
	control.Post("/delete", routerHandler.DeleteRoute)
	control.Post("/get", routerHandler.GetRoute)
	control.Post("/list", routerHandler.ListRoute)

}
