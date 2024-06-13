package message

type MessageIncoming struct {
	Src      string `form:"src"`
	Dst      string `form:"dst"`
	TaskId   string `form:"taskid"`
	MetaData string `form:"metadata"`
}

type Message struct {
	File     string `json:"file"`
	Length   int    `json:"length"`
	MetaData string `json:"metadata"`
}
