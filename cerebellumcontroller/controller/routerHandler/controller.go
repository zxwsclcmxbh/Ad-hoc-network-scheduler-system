package routerHandler

import (
	"iai/cerebellumController/dao"
	"iai/cerebellumController/dao/model"
	"iai/cerebellumController/dto/common"
	"iai/cerebellumController/dto/control"

	"iai/cerebellumController/service/routerService"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(c *fiber.Ctx) error {
	req := new(control.AddRouteReq)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: "invalid request"})
	}
	if rid, err := routerService.AddRoute(model.RouteItem{Taskid: req.TaskId, Src: req.Src, Svc: req.Svc, Dst: req.Dst, Node: req.Node}); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		return c.JSON(common.Response{Code: 0, Msg: "ok", Data: control.AddRouteResp{RouteId: rid}})
	}
}

func DeleteRoute(c *fiber.Ctx) error {
	req := new(control.DeleteRouteReq)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: "invalid request"})
	}
	routerService.DeleteRoute(req.RouteId)
	return c.JSON(common.Response{Code: 0, Msg: "ok"})
}

func ModifyRoute(c *fiber.Ctx) error {
	req := new(control.ModifyRouteReq)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: "invalid request"})
	}
	if err := routerService.ModifyRoute(*req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	}
	return c.JSON(common.Response{Code: 0, Msg: "ok"})
}

func ListRoute(c *fiber.Ctx) error {
	req := new(control.ListRouteReq)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: "invalid request"})
	}
	if result, err := dao.GetRouteList(req.TaskId); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		return c.JSON(common.Response{Code: 0, Msg: "ok", Data: control.ListRouteResp{Length: len(result), Routes: result}})
	}
}

func GetRoute(c *fiber.Ctx) error {
	req := new(control.GetRouteReq)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: "invalid request"})
	}
	if r, err := dao.GetRoute(req.RouteId); err != nil {
		return c.JSON(common.Response{Code: -1, Msg: err.Error()})
	} else {
		return c.JSON(common.Response{Code: 0, Msg: "ok", Data: control.GetRouteResp{r}})

	}
}
