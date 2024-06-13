package model

import (
	"github.com/jinzhu/gorm"
)

type RouteItem struct {
	gorm.Model
	RouteId string `json:"route_id" gorm:"index;size:32" mapstructure:"route_id"`
	Dst     string `json:"dst"`
	Src     string `json:"src"`
	Node    string `json:"node"`
	Taskid  string `json:"task_id" mapstructure:"task_id"`
	Svc     string `json:"svc"`
}
