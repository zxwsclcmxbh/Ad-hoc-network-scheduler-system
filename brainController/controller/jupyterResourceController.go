package controller

import (
	"bytes"
	"cloud/brainController/common"
	"cloud/brainController/service/dockerService"
	"cloud/brainController/service/resourceService"
	"cloud/brainController/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
)

type CreateJupyterResourceReq struct {
	ProjectName     string `json:"project_name"`
	DockerTarUrl    string `json:"docker_tar_url"`
	ImageName       string `json:"image_name"`
	ImageTagVersion string `json:"image_tag_version"`
	ImagePort       int32  `json:"image_port"`
	IsGPU           bool   `json:"is_gpu"`
}

// type CreateJupyterResourceResponse struct {
// 	common.BaseResponse
// 	Token string           `json:"token"`
// 	IP    string           `json:"IP"`
// 	Port  []v1.ServicePort `json:"port"`
// }

type CreateJupyterResourceResponse struct {
	common.BaseResponse
	Password string `json:"password"`
	URL      string `json:"url"`
}

type ProjectInfo struct {
	ProjectName string ` json:"project_name"`
}

type GetJupyterInfoOutput struct {
	Token string           `json:"token"`
	IP    string           `json:"IP"`
	Port  []v1.ServicePort `json:"port"`
}

type PushJupyterImageResponse struct {
	common.BaseResponse
	BuildCmdOut string `json:"build_cmd_out"`
	PullCmdOut  string `json:"pull_cmd_out"`
}

func CreateJupyterResource(c *gin.Context) {
	clientSet := utils.ClientSet
	var req CreateJupyterResourceReq
	if c.ShouldBind(&req) == nil {
		log.Println("params:")
		log.Println(req)
	}
	rid := utils.GetRandomString(20)
	password, hex := utils.GenPassword()
	log.Println(rid, password, hex)
	var nodeName string
	if req.IsGPU {
		nodeName = common.GPUNode
	} else {
		nodeName = "cloud-node"
	}
	log.Println("nodeName:" + nodeName)

	imageName := common.HubAddress + ":" + common.HubPort + "/" + req.ImageName + ":" + req.ImageTagVersion

	//build镜像
	buildInfo, buildImageErr := dockerService.BuildImage(req.DockerTarUrl, imageName)
	log.Println("BuildImage:" + buildInfo)
	if buildImageErr != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "error:" + buildImageErr.Error(),
		})
		return
	}

	//TODO:判断本地是否存在该镜像

	//push镜像
	pushInfo, pushImageErr := dockerService.PushImage(imageName)
	log.Println("pushImage:" + pushInfo)
	if pushImageErr != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "pushImageError:" + pushImageErr.Error(),
		})
		return
	}

	// 创建deployment
	flagDeploy, deployErr := resourceService.CreatePHMDeployment(clientSet, "phm", req.ProjectName+"-deploy", 1, req.ProjectName+"-label", nodeName, req.ProjectName+"-container", imageName, req.ImagePort, rid, password, hex, req.IsGPU)
	// flagDeploy, deployErr := resourceService.CreateDeployment(clientSet, "phm", req.ProjectName+"-deploy", 1, req.ProjectName+"-label", nodeName, req.ProjectName+"-container", imageName, req.ImagePort, req.IsGPU)
	if deployErr != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "deployError:" + deployErr.Error(),
		})
		return
	}
	// 创建service
	// flagSvc, svcErr := resourceService.CreateService(clientSet, "phm", req.ProjectName+"-svc", req.ImagePort, req.ProjectName+"-label")
	flagSvc, svcErr := resourceService.CreatePHMService(clientSet, "phm", req.ProjectName+"-svc", req.ImagePort, req.ProjectName+"-label")
	if svcErr != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "svcError:" + svcErr.Error(),
		})
		return
	}

	//判断是否部署成功
	if flagSvc && flagDeploy {
		if err := resourceService.WaitDeploymentReady(clientSet, "phm", req.ProjectName+"-label"); err != nil {
			resourceService.DeleteSvcAndDeploy(clientSet, "phm", req.ProjectName+"-deploy", req.ProjectName+"-svc")
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: 500,
				StatusMsg:  "faileds:CreateJupyterResourceTiemout",
			})
			return
		} else {
			time.Sleep(time.Second * 5)
			// jupyterInfo, _ := getJupyterInfo(req.ProjectName)
			log.Println(utils.Addroute(req.ProjectName+"-svc", rid, "mlflow"))
			log.Println(utils.Addroute(req.ProjectName+"-svc", rid, "jupyter"))
			c.JSON(http.StatusOK, CreateJupyterResourceResponse{
				BaseResponse: common.BaseResponse{
					StatusCode: 0,
					StatusMsg:  "success:CreateJupyterResource",
				},
				Password: password,
				URL:      utils.Config.Ingress.BaseUrl + rid + "/",
			})
			return
		}
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 2,
		StatusMsg:  "failed:CreateJupyterResource",
	})
}

func DeleteJupyterResource(c *gin.Context) {
	clientSet := utils.ClientSet
	var projectInfo ProjectInfo
	if c.ShouldBind(&projectInfo) == nil {
		log.Println(projectInfo)
	}
	fmt.Println("project-name:" + projectInfo.ProjectName)

	// 删除deployment和service
	flag, err := resourceService.DeleteSvcAndDeploy(clientSet, "phm", projectInfo.ProjectName+"-deploy", projectInfo.ProjectName+"-svc")
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "error:" + err.Error(),
		})
		return
	}

	//判断是否删除成功
	if flag {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success:DeleteJupyterResource",
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 2,
		StatusMsg:  "failed:DeleteJupyterResource",
	})
}

func getJupyterInfo(projectName string) (*GetJupyterInfoOutput, error) {
	clientSet := utils.ClientSet
	//根据projectName获取podName
	podsList := resourceService.GetPodsListInNamespace(clientSet, "phm")
	log.Println(podsList)
	podName := ""
	for _, item := range podsList {
		arr := strings.Split(item, "-deploy")
		if arr[0] == projectName {
			podName = item
			break
		}
	}

	//根据projectName获取svcName
	svcList := resourceService.GetSvcListInNamespace(clientSet, "phm")
	log.Println(svcList)
	svcName := ""
	for _, item := range svcList {
		arr := strings.Split(item, "-svc")
		if arr[0] == projectName {
			svcName = item
			break
		}
	}

	if podName == "" && svcName == "" {
		return nil, errors.New("projectName not found")
	}

	podInfo, getPodErr := resourceService.GetPodInfo(clientSet, "phm", podName)
	log.Println(podInfo)
	if getPodErr != nil {
		return nil, getPodErr
	}
	svcInfo, getSvcErr := resourceService.GetSvcInfo(clientSet, "phm", svcName)
	if getSvcErr != nil {
		return nil, getSvcErr
	}
	var token string
	var err error
	for count := 0; count <= 5; count++ {
		token, err = getToken(podName)
		if err != nil {
			continue
		} else {
			break
		}
	}
	return &GetJupyterInfoOutput{Token: token,
		IP:   common.NodeIP[podInfo.Spec.NodeName],
		Port: svcInfo.Spec.Ports}, nil
}

func getToken(podName string) (string, error) {
	name, getTokenErr := resourceService.WatchContainerLogWithPodName(utils.ClientSet, "phm", podName)
	if getTokenErr != nil {
		return "", getTokenErr
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(name)
	newName := buf.String()
	if strings.Contains(newName, "Use Control-C to stop") {
		return newName, nil
	} else {
		return "", errors.New("token not found")
	}
}
