package cerebellumService

import (
	"cloud/brainController/model"
	"cloud/brainController/utils"
	"encoding/json"
	"log"

	"github.com/olahol/melody"
)

func SendMessageToCerebellum(msg interface{}, nodename string) error {
	msg_byte, _ := json.Marshal(msg)
	log.Println(string(msg_byte))
	sessions, _ := utils.M.Sessions()
	sess_to_send := []*melody.Session{}
	for _, sess := range sessions {
		if val, _ := sess.Get("nodename"); val.(string) == nodename {
			sess_to_send = append(sess_to_send, sess)
		}
	}
	if len(sess_to_send) == 0 {
		log.Println("node not online")
		if cache, ok := utils.MessageCache.Load(nodename); ok {
			c := cache.(*[][]byte)
			*c = append(*c, msg_byte)
			log.Println("save messgae to cache")
		} else {
			c := [][]byte{msg_byte}
			utils.MessageCache.Store(nodename, &c)
			log.Println("make new cache for node", nodename)
		}
		return nil
	} else {
		return utils.M.BroadcastMultiple(msg_byte, sess_to_send)
	}
	// return utils.M.BroadcastFilter(msg_byte, func(s *melody.Session) bool {
	// 	if val, _ := s.Get("nodename"); val.(string) == nodename {
	// 		return true
	// 	}
	// 	return false
	// })
}

func SendRoute(msg model.AddRouteReq, nodename string) error {
	msg_full := model.CerebellumMessage{Src: "brain", Dst: nodename, Type: "route:add", Payload: msg}
	return SendMessageToCerebellum(msg_full, nodename)
}

func DeleteRoute(msg model.DeleteRouteReq, nodename string) error {
	msg_full := model.CerebellumMessage{Src: "brain", Dst: nodename, Type: "route:delete", Payload: msg}
	return SendMessageToCerebellum(msg_full, nodename)
}

func ModifyRoute(msg model.ModifyRouteReq, nodename string) error {
	msg_full := model.CerebellumMessage{Src: "brain", Dst: nodename, Type: "route:modify", Payload: msg}
	return SendMessageToCerebellum(msg_full, nodename)
}

func SendMessage(cerebellumch chan model.MessageWithNode) {
	for i := range cerebellumch {
		switch msg := i.Msg.(type) {
		case model.AddRouteReq:
			SendRoute(msg, i.Node)
		case model.DeleteRouteReq:
			DeleteRoute(msg, i.Node)
		default:
			log.Println("no such type")
		}
	}
}
