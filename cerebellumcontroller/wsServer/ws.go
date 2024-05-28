package wsServer

import (
	"iai/cerebellumController/config"
	"iai/cerebellumController/dao/model"
	"iai/cerebellumController/dto/control"
	"iai/cerebellumController/dto/ws"
	"iai/cerebellumController/service/brain"
	"iai/cerebellumController/service/routerService"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var C *websocket.Conn

func Run(interrupt chan os.Signal) {
	var err error
	C, _, err = websocket.DefaultDialer.Dial(config.Config.Brain, nil)
	if err != nil {
		// log.Fatalf("cannot connect to %s for %v", config.Config.Brain, err.Error())
		if config.Config.Debug == "false" {
			log.Fatalf("cannot connect to %s for %v", config.Config.Brain, err.Error())
		} else {
			return
		}

	}
	defer C.Close()
	login := ws.LoginPayload{NodeName: config.Config.Host}
	if err := brain.SendLogin(login, C); err != nil {
		if config.Config.Debug == "false" {
			log.Fatal("login unsuccess:", err)
		}
	}
	recv_channel := make(chan ws.CerebellumMessage)
	close_channel := make(chan interface{}, 1)
	go func() {
		defer close(close_channel)
		for {
			m := ws.CerebellumMessage{}
			if err := C.ReadJSON(&m); err != nil {
				log.Fatalln("recv msg with error:", err)
				return
			} else {
				recv_channel <- m
			}
		}
	}()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-close_channel:
			return
		case <-ticker.C:
			if err := C.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				log.Fatalln("cannont ping brain with error:", err)
			}
		case <-interrupt:
			log.Println("stop")
			if err := C.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(5*time.Second)); err != nil {
				return
			}
			select {
			case <-close_channel:
			case <-time.After(time.Second):
			}
			return
		case msg := <-recv_channel:
			log.Println(msg)
			switch msg.Type {
			case "route:add":
				r := &control.AddRouteReq{}
				mapstructure.Decode(msg.Payload, r)
				log.Println(routerService.AddRoute(model.RouteItem{Taskid: r.TaskId, Src: r.Src, Svc: r.Svc, Dst: r.Dst, Node: r.Node}))
			case "route:modify":
				r := &control.ModifyRouteReq{}
				r.RouteId = r.Dst
				mapstructure.Decode(msg.Payload, r)
				log.Println(routerService.ModifyRoute(*r))
			case "route:delete":
				r := &control.DeleteRouteReq{}
				mapstructure.Decode(msg.Payload, r)
				log.Println(routerService.DeleteRoute(r.RouteId))
			}
		}
	}
}
