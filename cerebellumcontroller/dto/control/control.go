package control

import (
	"iai/cerebellumController/dao/model"
)

type AddRouteReq struct {
	TaskId string `json:"task_id"  mapstructure:"task_id"` // task_id 全局任务编号
	Dst    string `json:"dst"`                             // 目标队列名称
	Node   string `json:"node"`                            // 节点名称
	Svc    string `json:"svc"`                             // 节点内redis服务名称
	Src    string `json:"src"`
}

type AddRouteResp struct {
	RouteId string `json:"route_id"` // 路由条目id
}

type DeleteRouteReq struct {
	RouteId string `json:"route_id"  mapstructure:"route_id"` //删除路由条目id
}

type ModifyRouteReq struct {
	model.RouteItem
}

type ListRouteReq struct {
	TaskId string `json:"task_id"` // 要列出的TaskId
}

type ListRouteResp struct {
	Length int               `json:"length"`
	Routes []model.RouteItem `json:"routes"`
}

type GetRouteReq struct {
	RouteId string `json:"route_id"`
}

type GetRouteResp struct {
	model.RouteItem
}
