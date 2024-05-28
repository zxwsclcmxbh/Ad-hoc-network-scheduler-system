package model

import "gorm.io/gorm"

type TaskDefinition struct {
	Image      string `json:"image"`
	Node       string `json:"node"`
	CustomConf string `json:"custom_conf"`
	Type       string `json:"type"`
}

type Task struct {
	gorm.Model
	TaskCreateRequest
	TaskStatus string
	TaskId     string        `gorm:"index;size:32"`
	IntakePort []int         `gorm:"serializer:json"`
	UserId     string        `gorm:"index;size:32"`
	Pods       []*Pod        `gorm:"foreignKey:TaskId;references:TaskId"`
	Queues     []*QueueRoute `gorm:"foreignKey:TaskId;references:TaskId"`
}

type TaskOld struct {
	gorm.Model
	TaskCreateRequestOld
	TaskStatus string
	TaskId     string `gorm:"index;size:32"`
	IntakePort int
	UserId     string
	Pods       []*Pod        `gorm:"foreignKey:TaskId;references:TaskId"`
	Queues     []*QueueRoute `gorm:"foreignKey:TaskId;references:TaskId"`
}

type Pod struct {
	gorm.Model
	PodId        string `gorm:"index;size:32"`
	TaskId       string `gorm:"index;size:32"`
	Image        string `gorm:"size:512"`
	Node         string
	Type         string
	Environment  map[string]string `gorm:"serializer:json"`
	NodeAffinity string
}

type QueueRoute struct {
	gorm.Model
	QueueId  string `gorm:"index;size:32"`
	TaskId   string `gorm:"index;size:32"`
	SrcPod   string
	DstPod   string
	IsRemote bool
	SrcNode  string
	DstNode  string
}

type IaiImage struct {
	gorm.Model
	ImageId       string `gorm:"index;size:32"`
	UserId        string `gorm:"index;size:32"`
	ImageNickName string //镜像中文名称
	Image         string `gorm:"size:512"` //镜像名称
	ImageDes      string //镜像描述
	ImagePort     int32
	IsPublic      *bool `gorm:"default:true"`
}

type Node struct {
	Id   string `json:"id"`
	Meta struct {
		Id string `json:"id"`
	} `json:"meta"`
}
type Link struct {
	Id      string `json:"id"`
	StartId string `json:"startid"`
	EndId   string `json:"endid"`
}
type Graph struct {
	NodeList []Node `json:"nodelist"`
	LinkList []Link `json:"linklist"`
}

type TaskDefination struct {
	Graph Graph       `json:"graph"`
	Data  interface{} `json:"data"`
}
