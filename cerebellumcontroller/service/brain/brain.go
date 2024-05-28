package brain

import (
	"encoding/json"
	"iai/cerebellumController/config"
	"iai/cerebellumController/dto/ws"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Conn    *websocket.Conn
	Message []byte
}

var MessageChan chan Message

func SendWorker() {
	MessageChan = make(chan Message)
	for mess := range MessageChan {
		mess.Conn.WriteMessage(websocket.TextMessage, mess.Message)
	}
}

func Upstream(msg interface{}, conn *websocket.Conn) error {
	text, _ := json.Marshal(msg)
	log.Println(string(text))
	m := Message{conn, text}
	MessageChan <- m
	return nil
}

func SendLogin(l ws.LoginPayload, conn *websocket.Conn) error {
	m := ws.CerebellumMessage{Src: config.Config.Host, Dst: "brain", Type: "login", Payload: l}
	return Upstream(&m, conn)
}

func SendLog(l ws.ReportLogPayload, conn *websocket.Conn) error {
	m := ws.CerebellumMessage{Src: config.Config.Host, Dst: "brain", Type: "report:podlog", Payload: l}
	log.Println(m)
	return Upstream(&m, conn)
}

func SendValue(l ws.ReportMessagePayload, conn *websocket.Conn) error {
	log.Println(l)
	m := ws.CerebellumMessage{Src: config.Config.Host, Dst: "brain", Type: "report:value", Payload: l}
	return Upstream(&m, conn)
}

func SendTrace(l ws.ReportTracePayload, conn *websocket.Conn) error {
	m := ws.CerebellumMessage{Src: config.Config.Host, Dst: "brain", Type: "report:trace", Payload: l}
	return Upstream(&m, conn)
}
