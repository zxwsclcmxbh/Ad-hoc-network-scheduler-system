package model

type TaskCreateRequestOld struct {
	TaskName       string           `json:"name"`
	TaskDescribe   string           `json:"description"`
	TaskDefinition []TaskDefinition `gorm:"serializer:json" json:"definition"`
}

type TaskCreateRequest struct {
	TaskName       string `json:"name"`
	TaskDescribe   string `json:"description"`
	TaskDefinition string `json:"definition"`
}

type CerebellumMessage struct {
	Src     string      `json:"src"`
	Dst     string      `json:"dst"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type LoginPayload struct {
	NodeName string `json:"nodename"`
}

type AddRouteReq struct {
	TaskId string `json:"task_id"` // task_id 全局任务编号
	Dst    string `json:"dst"`     // 目标队列名称
	Node   string `json:"node"`    // 节点名称
	Svc    string `json:"svc"`     // 目的节点内redis服务名称
	Src    string `json:"src"`     //信源队列名称
}

type DeleteRouteReq struct {
	RouteId string `json:"route_id"` //删除路由条目id
}

type ModifyRouteReq struct {
	RouteId string `json:"route_id" gorm:"primaryKey"`
	Dst     string `json:"dst"`
	Src     string `json:"src"`
	Node    string `json:"node"`
	Taskid  string `json:"task_id"`
	Svc     string `json:"svc"`
}

type ReportLogPayload struct {
	TaskId    string                 `json:"task_id" mapstructure:"task_id"`
	PodName   string                 `json:"pod_name" mapstructure:"pod_name"`
	TimeStamp string                 `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

type ReportMessagePayload struct {
	EquipmentId  string              `json:"equipment_id" mapstructure:"equipment_id"`
	TaskId       string              `json:"task_id" mapstructure:"task_id"`
	MessageItems []ReportMessageItem `json:"items" mapstructure:"items"`
}

type ReportMessageItem struct {
	TimeStamp string                 `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

type MessageWithNode struct {
	Msg  interface{}
	Node string
}

type ReportTracePayload struct {
	TaskId    string `json:"task" mapstructure:"task"`
	PodId     string `json:"pod" mapstructure:"pod"`
	TimeStamp int64  `json:"time" mapstructure:"time"`
	TraceType string `json:"type" mapstructure:"type"`
	MessageId string `json:"message_id" mapstructure:"message_id"`
}
