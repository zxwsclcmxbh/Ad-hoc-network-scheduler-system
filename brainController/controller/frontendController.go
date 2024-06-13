package controller

import (
	"cloud/brainController/common"
	"cloud/brainController/service/nodeService"
	"cloud/brainController/service/resourceService"
	"cloud/brainController/service/taskService"
	"cloud/brainController/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseCountResponse struct {
	common.BaseResponse
	Data struct {
		NodeCount int
		PodCount  int
		NSCount   int
		TaskCount int
	} `json:"data"`
}

func GetBaseCount(c *gin.Context) {
	nodelist, _ := nodeService.GetNodeList(utils.ClientSet)
	namespacelist := resourceService.GetNamespaces(utils.ClientSet)
	var podList []string
	for _, item := range namespacelist {
		t := resourceService.GetPodsListInNamespace(utils.ClientSet, item)
		podList = append(podList, t...)
	}
	taskslist := taskService.GetTaskList()
	c.JSON(http.StatusOK, BaseCountResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Data: struct {
			NodeCount int
			PodCount  int
			NSCount   int
			TaskCount int
		}{
			NodeCount: len(nodelist),
			PodCount:  len(podList),
			TaskCount: len(taskslist),
			NSCount:   len(namespacelist),
		},
	})
}

type UserInfo struct {
	Code int `json:"code"`
	Data struct {
		Roles        []string `json:"roles"`
		Name         string   `json:"name"`
		Avatar       string   `json:"avatar"`
		Introduction string   `json:"introduction"`
	} `json:"data"`
}

func GetUserDetail(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, UserInfo{
		Code: 20000,
		Data: struct {
			Roles        []string `json:"roles"`
			Name         string   `json:"name"`
			Avatar       string   `json:"avatar"`
			Introduction string   `json:"introduction"`
		}{
			Roles:        []string{"admin"},
			Introduction: "云边协同",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Name:         username,
		},
	})
}
