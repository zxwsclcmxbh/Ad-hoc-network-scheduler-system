package controller

import (
	"bytes"
	"cloud/brainController/model"
	"cloud/brainController/utils"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/olahol/melody"
)

func HandleWSReq(c *gin.Context) {
	utils.M.HandleRequest(c.Writer, c.Request)

}

func HandlerMessage(s *melody.Session, msg []byte) {
	msg_struct := &model.CerebellumMessage{}
	dec := json.NewDecoder(bytes.NewReader(msg))
	dec.UseNumber()
	if err := dec.Decode(msg_struct); err != nil {
		log.Println("msg header decode error")
	}
	switch msg_struct.Type {
	case "login":
		HandleLogin(msg_struct.Payload, s)
	case "report:podlog":
		HandleLog(msg_struct.Payload, s)
	case "report:value":
		go HandleReport(msg_struct.Payload, s)
	case "report:trace":
		go HandleTrace(msg_struct.Payload, s)
	}

}

func HandleTrace(msg interface{}, s *melody.Session) {
	l := model.ReportTracePayload{}
	// log.Println(l)
	mapstructure.Decode(msg, &l)
	t := time.Unix(0, l.TimeStamp)
	// log.Println("time convert", t)
	p, _ := client.NewPoint("trace", map[string]string{"message_id": l.MessageId, "task_id": l.TaskId, "src_node": s.MustGet("nodename").(string), "pod_id": l.PodId}, map[string]interface{}{"trace_type": l.TraceType}, t)
	// log.Println("add new point", p, err)
	utils.WriteCh <- p
}

func HandleLogin(msg interface{}, s *melody.Session) {
	l := model.LoginPayload{}
	mapstructure.Decode(msg, &l)
	log.Printf("%s logged in\n", l.NodeName)
	s.Set("nodename", l.NodeName)
	if val, ok := utils.MessageCache.Load(l.NodeName); ok {
		c := val.(*[][]byte)
		if len(*c) != 0 {
			for _, msg := range *c {
				s.Write(msg)
				log.Println("sync messgae to node", l.NodeName, "with messages", string(msg))
			}
			utils.MessageCache.Delete(l.NodeName)
		}
	}
}

func HandleReport(msg interface{}, s *melody.Session) {
	r := model.ReportMessagePayload{}
	mapstructure.Decode(msg, &r)
	log.Println("receive report")
	for _, item := range r.MessageItems {
		// log.Println(item.TimeStamp)
		ts, _ := strconv.ParseFloat(item.TimeStamp, 64)
		second := int(math.Floor(float64(ts)))
		nanosecond := int((ts - float64(second)) * 1e9)
		// log.Println(second, nanosecond)
		t := time.Unix(int64(second), int64(nanosecond))
		// log.Println("time convert", t)
		p, _ := client.NewPoint("value", map[string]string{"equipment_id": r.EquipmentId, "task_id": r.TaskId, "src_node": s.MustGet("nodename").(string)}, item.Values, t)
		// log.Println("add new point", p, err)
		utils.WriteCh <- p
	}

	// TODO: 记录入数据库
}

func HandleLog(msg interface{}, s *melody.Session) {
	r := model.ReportLogPayload{}
	mapstructure.Decode(msg, &r)
	log.Println("receive log")
	ts, _ := strconv.ParseFloat(r.TimeStamp, 64)
	second := int(math.Floor(float64(ts)))
	nanosecond := int((ts - float64(second)) * 1e9)
	// log.Println(second, nanosecond)
	// t := time.Unix(int64(second), int64(nanosecond))
	// log.Println("time convert", t)
	p, _ := client.NewPoint("log", map[string]string{"podname": r.PodName, "task_id": r.TaskId, "src_node": s.MustGet("nodename").(string)}, r.Values, time.Unix(int64(second), int64(nanosecond)))
	utils.WriteCh <- p
	// TODO: 记录入数据库
}
