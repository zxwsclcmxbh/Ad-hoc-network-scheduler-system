package routerService

import (
	"context"
	"errors"
	"iai/cerebellumController/dao"
	"iai/cerebellumController/dao/model"
	"iai/cerebellumController/dto/control"
	"iai/cerebellumController/service/worker"
	"sync"
)

var routeMap sync.Map

func Init() {
	list, _ := dao.GetAllRouteList()
	for _, item := range list {
		ctx, c := context.WithCancel(context.Background())
		go worker.HandleOutgoingMessage(item.Taskid, item.Src, item.Node, item.Dst, item.Svc, ctx)
		routeMap.Store(item.RouteId, c)
	}
}
func AddRoute(r model.RouteItem) (string, error) {
	r.RouteId = r.Dst
	if err := dao.AddRoute(r); err != nil {
		return "", err
	} else {
		ctx, c := context.WithCancel(context.Background())
		go worker.HandleOutgoingMessage(r.Taskid, r.Src, r.Node, r.Dst, r.Svc, ctx)
		routeMap.Store(r.RouteId, c)
		return r.RouteId, nil
	}
}
func DeleteRoute(routeid string) error {
	if c, ok := routeMap.Load(routeid); ok {
		if cancel, ok := c.(context.CancelFunc); ok {
			cancel()
			dao.DeleteRoute(routeid)
			return nil
		} else {
			return errors.New("unknown error")
		}
	} else {
		return errors.New("routeid not exists")
	}
}
func ModifyRoute(r control.ModifyRouteReq) error {
	if c, ok := routeMap.Load(r.RouteId); ok {
		if cancel, ok := c.(context.CancelFunc); ok {
			cancel()
			dao.ModifyRoute(r.RouteItem)
			ctx, c := context.WithCancel(context.Background())
			go worker.HandleOutgoingMessage(r.Taskid, r.Src, r.Node, r.Dst, r.Svc, ctx)
			routeMap.Store(r.RouteId, c)
			return nil
		} else {
			return errors.New("unknown error")
		}
	} else {
		return errors.New("routeid not exists")
	}
}
