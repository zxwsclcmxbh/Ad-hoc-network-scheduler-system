package utils

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
)

type Nodes struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}
type Upstream struct {
	Nodes   []Nodes `json:"nodes"`
	Timeout struct {
		Connect int `json:"connect"`
		Send    int `json:"send"`
		Read    int `json:"read"`
	} `json:"timeout"`
	Type          string `json:"type"`
	Scheme        string `json:"scheme"`
	PassHost      string `json:"pass_host"`
	KeepAlivePool struct {
		IdleTimeout int `json:"idle_timeout"`
		Requests    int `json:"requests"`
		Size        int `json:"size"`
	} `json:"keepalive_pool"`
}
type AddRouteReq struct {
	// ID              string   `json:"id"`
	Uri             string      `json:"uri"`
	Name            string      `json:"name"`
	Methods         []string    `json:"methods"`
	Upstream        Upstream    `json:"upstream"`
	EnableWebsocket bool        `json:"enable_websocket"`
	Status          int         `json:"status"`
	Plugins         interface{} `json:"plugins,omitempty"`
}
type Resp struct {
	Code    int
	Message string
}

var rtypemap = map[string]int{
	"jupyter": 8888,
	"mlflow":  8889,
}

func Addroute(targetname string, rid string, rtype string) error {
	req := AddRouteReq{
		// ID:      targetname + "-" + rtype,
		Uri:     "/" + rid + "/" + rtype + "*",
		Name:    targetname + "-" + rtype,
		Methods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"},
		Upstream: Upstream{
			Nodes: []Nodes{
				{
					Host:   targetname,
					Port:   rtypemap[rtype],
					Weight: 1,
				},
			},
			Timeout: struct {
				Connect int "json:\"connect\""
				Send    int "json:\"send\""
				Read    int "json:\"read\""
			}{6, 1000, 1000},
			Scheme:   "http",
			PassHost: "pass",
			KeepAlivePool: struct {
				IdleTimeout int "json:\"idle_timeout\""
				Requests    int "json:\"requests\""
				Size        int "json:\"size\""
			}{1000, 1000, 320},
		},
		EnableWebsocket: true,
		Status:          1,
	}
	if rtype == "mlflow" {
		config := make(map[string][]string)
		config["regex_uri"] = []string{"^/" + rid + "/mlflow(.*)", "$1"}
		plugin := make(map[string]interface{})
		plugin["proxy-rewrite"] = config
		req.Plugins = plugin
	}
	resp, err := resty.New().R().SetHeader("X-API-KEY", Config.Ingress.ApiKey).SetBody(req).Post(Config.Ingress.AdminUrl + "apisix/admin/routes")
	body := resp.Body()
	log.Println(string(body))
	if err != nil {
		return err
	} else {
		if resp.StatusCode() != 200 {
			res := Resp{}
			if err := json.Unmarshal(body, &res); err != nil {
				return err
			} else {
				return errors.New(res.Message)
			}
		} else {
			return nil
		}
	}
}
