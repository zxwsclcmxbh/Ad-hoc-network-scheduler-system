package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iai/cerebellumController/config"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GenFileName(basepath string, originname string, taskid string) (newpath string) {
	u := uuid.New()
	fname := u.String()
	ext := filepath.Ext(originname)
	nfname := fmt.Sprintf("%s%s", fname, ext)
	newpath = filepath.Join(basepath, taskid, nfname)
	return
}

func GenRouteId() (routeid string) {
	u := uuid.New()
	routeid = u.String()
	return
}

func GenHost(svc string, node string) (host string) {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:3000/", svc, "iai")
}

func ReGenFileNameOutgoing(orgin string, taskid string) string {
	if config.Config.Debug == "false" {
		return strings.Replace(orgin, "/data/", "/data/"+taskid+"/", 1)
	} else {
		return orgin
	}
}

func ReGenFileNameIncoming(orgin string, taskid string) string {
	if config.Config.Debug == "false" {
		return strings.Replace(orgin, "/data/"+taskid+"/", "/data/", 1)
	} else {
		return orgin
	}
}

func Struct2Map(content interface{}) map[string]interface{} {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		log.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
			log.Println(err)
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name
}
