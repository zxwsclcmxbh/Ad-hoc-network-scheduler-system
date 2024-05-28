package ws

type CerebellumMessage struct {
	Src     string      `json:"src"`
	Dst     string      `json:"dst"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type LoginPayload struct {
	NodeName string `json:"nodename"`
}

type ReportMessagePayload struct {
	EquipmentId  string              `json:"equipment_id"`
	TaskId       string              `json:"task_id"`
	MessageItems []ReportMessageItem `json:"items"`
}

type ReportMessageItem struct {
	TimeStamp string                 `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

type ReportLogPayload struct {
	TaskId    string                 `json:"task_id"`
	PodName   string                 `json:"pod_name"`
	TimeStamp string                 `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

type IncomingLog struct {
	Name            string `form:"name"`
	Msg             string `form:"msg"`
	Args            string `form:"args"`
	LevelName       string `form:"levelname"`
	LevelNo         int    `form:"levelno"`
	PathName        string `form:"pathname"`
	FileName        string `from:"filename"`
	Module          string `form:"module"`
	ExcInfo         string `form:"exc_info"`
	ExcText         string `form:"exc_text"`
	StatckInfo      string `form:"stack_info"`
	Lineno          int    `form:"lineno"`
	FuncName        string `form:"funcName"`
	Created         string `form:"created"`
	Msecs           string `form:"msecs"`
	RelativeCreated string `form:"relativeCreated"`
	Thread          string `form:"thread"`
	ThreadName      string `form:"threadName"`
	ProcessName     string `form:"processName"`
	Process         string `form:"process"`
}

type ReportTracePayload struct {
	TaskId    string `json:"task"`
	PodId     string `json:"pod"`
	TimeStamp int64  `json:"time"`
	TraceType string `json:"type"`
	MessageId string `json:"message_id"`
}
