package controller

import (
	"cloud/brainController/common"
	"cloud/brainController/service/nodeService"
	"cloud/brainController/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetNodeListResponse struct {
	common.BaseResponse
	NodeList []string `json:"node_list"`
}

type GetNodeListMetricsResponse struct {
	common.BaseResponse
	Data []nodeService.NodeInfoNew `json:"data"`
}

type GetNodeInfoResponse struct {
	common.BaseResponse
	NodeInfo map[string]interface{} `json:"node_info"`
}

func GetNodeList(c *gin.Context) {
	clientSet := utils.ClientSet
	nodeList, err := nodeService.GetNodeList(clientSet)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	var nodeListArray []string
	for _, nds := range nodeList {
		nodeListArray = append(nodeListArray, nds.Name)
	}
	c.JSON(http.StatusOK, GetNodeListResponse{
		NodeList: nodeListArray,
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success:GetNodeList",
		},
	})
}

func GetNodeListWithMetrics(c *gin.Context) {
	clientSet := utils.MetricsClientSet
	nodeList, err := nodeService.GetNodeListWithMetrics(clientSet, utils.ClientSet)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetNodeListMetricsResponse{
		Data: nodeList,
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success:GetNodeList",
		},
	})
}

func GetNodeInfo(c *gin.Context) {
	clientSet := utils.ClientSet
	nodeName := c.Query("node_name")
	nodeInfo, err := nodeService.GetNodeInfo(clientSet, nodeName)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetNodeInfoResponse{
		NodeInfo: nodeInfo,
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success:GetNodeInfo",
		},
	})

}
