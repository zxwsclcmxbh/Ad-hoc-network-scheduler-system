package router

import (
	"iai/cerebellumController/controller/reportHandler"

	"github.com/gofiber/fiber/v2"
)

func InitReport(app *fiber.App) {
	msg := app.Group("/api").Group("/v1").Group("/report")
	msg.Post("/log", reportHandler.ReportLog)
	msg.Post("/value", reportHandler.ReportValue)
	msg.Post("/trace", reportHandler.ReortTrace)
}
