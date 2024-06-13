package reportHandler

import (
	"bytes"
	"encoding/json"
	"iai/cerebellumController/dto/common"
	"iai/cerebellumController/dto/ws"
	"iai/cerebellumController/service/brain"
	"iai/cerebellumController/util"
	"iai/cerebellumController/wsServer"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ReportLog(c *fiber.Ctx) error {
	r := ws.IncomingLog{}
	if err := c.BodyParser(&r); err != nil {
		log.Println(err)
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		log.Println(r)
		t := strings.Split(r.Name, "/")
		podname := t[0]
		taskid := t[1]
		val := util.Struct2Map(r)
		l := ws.ReportLogPayload{TaskId: taskid, PodName: podname, TimeStamp: r.Created, Values: val}
		log.Println(brain.SendLog(l, wsServer.C))
		return c.JSON(common.Response{Code: 0, Msg: "ok"})
	}

}

func ReportValue(c *fiber.Ctx) error {
	var req ws.ReportMessagePayload
	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.UseNumber()
	if err := dec.Decode(&req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		brain.SendValue(req, wsServer.C)
		return c.JSON(common.Response{Code: 0, Msg: "ok"})
	}

}

func ReortTrace(c *fiber.Ctx) error {
	var req ws.ReportTracePayload
	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		log.Println(req)
		brain.SendTrace(req, wsServer.C)
		return c.JSON(common.Response{Code: 0, Msg: "ok"})
	}
	// return nil
}
