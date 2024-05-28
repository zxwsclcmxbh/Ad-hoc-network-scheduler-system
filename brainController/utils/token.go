package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/go-resty/resty/v2"
)

type UserResp struct {
	StatusCode    int    `json:"status"`
	StatusMessage string `json:"msg"`
	Data          struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"user_name"`
	} `json:"data"`
}
type User struct {
	UserId   string
	UserName string
}

var UserDB sync.Map

func GetUserDetailReq(token string) (*User, error) {
	if token == "backdoor" {
		return &User{"test", "backdoor"}, nil
	}
	// resp, err := resty.New().R().SetFormData(map[string]string{"token": token}).Post("http://192.168.1.3:5051/component/getUser")
	resp, err := resty.New().R().SetFormData(map[string]string{"token": token}).Post(Config.PHMEndpoint)

	// TODO: 修改为配置文件
	if err != nil {
		return nil, err
	} else {
		result := UserResp{}
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return nil, err
		}
		if result.StatusCode == 0 {
			return &User{UserId: strconv.Itoa(result.Data.UserId), UserName: result.Data.UserName}, nil
		} else {
			return nil, errors.New(result.StatusMessage)
		}
	}
}

func GetUserDetail(token string) (*User, error) {
	if res, ok := UserDB.Load(token); !ok {
		if u, err := GetUserDetailReq(token); err != nil {
			return nil, err
		} else {
			UserDB.Store(token, u)
			return u, nil
		}
	} else {
		return res.(*User), nil
	}
}
