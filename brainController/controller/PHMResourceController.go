package controller

import (
	"cloud/brainController/common"
	"cloud/brainController/model"
	"cloud/brainController/service/ossService"
	"cloud/brainController/service/resourceService"
	"cloud/brainController/service/taskService"
	"cloud/brainController/utils"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var getFile = []string{"ls", "-l", "/home/jovyan/work", "--time-style", "+%Y-%m-%dT%H:%M:%S%Z"}

type PushPHMResourceParams struct {
	ImageNickName   string `json:"image_nickName"`
	ImageDes        string `json:"image_des"`
	DockerTarUrl    string `json:"docker_tar_url"`
	ImageName       string ` json:"image_name"`
	ImageTagVersion string `json:"image_tag_version"`
	ImagePort       int32  ` json:"image_port"`
	UserId          string `json:"user_id"`
}

func PushPHMResource(c *gin.Context) {
	var db = utils.DB
	var req PushPHMResourceParams
	if c.ShouldBind(&req) != nil {
		log.Println("params:")
		log.Println(req)
	}
	imageName := common.HubAddress + ":" + common.HubPort + "/" + req.ImageName + ":" + req.ImageTagVersion
	log.Println("imageName:" + imageName)

	//build镜像
	// buildInfo, buildImageErr := dockerService.BuildImage(req.DockerTarUrl, imageName)
	// log.Println("BuildImage:" + buildInfo)
	// if buildImageErr != nil {
	// 	c.JSON(http.StatusOK, common.BaseResponse{
	// 		StatusCode: 1,
	// 		StatusMsg:  "error:" + buildImageErr.Error(),
	// 	})
	// 	return
	// }

	// //push镜像
	// pushInfo, pushImageErr := dockerService.PushImage(imageName)
	// log.Println("pushImage:" + pushInfo)
	// if pushImageErr != nil {
	// 	c.JSON(http.StatusOK, common.BaseResponse{
	// 		StatusCode: 1,
	// 		StatusMsg:  "pushImageError:" + pushImageErr.Error(),
	// 	})
	// 	return
	// }
	if err := utils.RunDockerBuildX(req.DockerTarUrl, imageName); err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "pushImageError:" + err.Error(),
		})
		return
	}
	//存信息入库
	iaiImage := model.IaiImage{
		ImageId:       taskService.GenUUID(),
		ImageNickName: req.ImageNickName,
		Image:         imageName,
		ImageDes:      req.ImageDes,
		ImagePort:     req.ImagePort,
		UserId:        req.UserId,
	}
	err := db.Create(&iaiImage).Error
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "createIaiImageError:" + err.Error(),
		})
		return
	}
	db.Save(&iaiImage)

	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success:PushPHMResource",
	})
}

func DeletePHMResource(c *gin.Context) {
	var db = utils.DB
	var IaiImage = model.IaiImage{}
	image := c.Query("image")
	log.Println("image:" + image)
	if image == "" {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 2,
			StatusMsg:  "error:image参数为空",
		})
	}

	//删除数据库中的record
	res := db.Where("image=?", image).Delete(&IaiImage)
	if res.Error == nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success:PushPHMResource",
		})
	}

	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 1,
		StatusMsg:  "error:" + res.Error.Error(),
	})
}

type SearchResultItem struct {
	ImageId       string `json:"id"`
	ImageNickName string `json:"title"`
	ImageDes      string `json:"desc"`
	Image         string `json:"image"`
	IsPublic      string `json:"isPublic"`
	UserId        string `json:"userId"`
}
type SearchResultResponse struct {
	common.BaseResponse
	Data []SearchResultItem `json:"data"`
}

func SearchImage(c *gin.Context) {
	q := c.Query("search")
	uid := c.GetString("uid")
	var db = utils.DB
	var result []SearchResultItem
	db.Table("iai_images").Where(db.Where("(image_nick_name like ? or image_des like ?) and (user_id=? or is_public=true)", "%"+q+"%", "%"+q+"%", uid)).Find(&result)
	c.JSON(http.StatusOK, SearchResultResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Data: result,
	})
}

func GetIaiImageList(c *gin.Context) {
	var db = utils.DB
	isTotal := c.Query("isTotal")
	var result []SearchResultItem
	uid := c.GetString("uid")
	log.Println("uid:" + uid)
	if isTotal == "true" {
		db.Table("iai_images").Where("is_public=?", 1).Find(&result)
	} else {
		db.Table("iai_images").Where("user_id=?", uid).Find(&result)
	}
	c.JSON(http.StatusOK, SearchResultResponse{
		BaseResponse: common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Data: result,
	})
}

type EditIaiImageRecordParams struct {
	ImageId       string `json:"id"`
	ImageNickName string `json:"title"`
	ImageDes      string `json:"desc"`
	IsPublic      bool   `json:"isPublic"`
}

func EditIaiImageRecord(c *gin.Context) {
	var db = utils.DB
	var req EditIaiImageRecordParams
	if c.ShouldBind(&req) != nil {
		log.Println("params:")
		log.Println(req)
	}

	res := db.Table("iai_images").Where("image_id=?", req.ImageId).Updates(model.IaiImage{IsPublic: &req.IsPublic, ImageNickName: req.ImageNickName, ImageDes: req.ImageDes})
	if res.Error == nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		})
		return
	}
	log.Println(res.Error)
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 1,
		StatusMsg:  "failed:EditIaiImageRecord",
	})

}

type GetFileListReq struct {
	ProjectName   string ` json:"project_name"`
	ContainerName string `json:"container_name"`
}
type FileItem struct {
	Name         string     `json:"name"`
	Size         int        `json:"size"`
	ModifiedTime *time.Time `json:"modifed_time"`
}

type GetFileReq struct {
	FileName      string `json:"file_name"`
	ProjectName   string `json:"project_name"`
	ContainerName string `json:"container_name"`
}
type OssFileResp struct {
	OriginFileName string `json:"origin_file_name"`
	URL            string `json:"url"`
}

func GetPodNameFromProject(project string, clientSet *kubernetes.Clientset) (string, error) {
	//根据projectName获取podName
	podsList := resourceService.GetPodsListInNamespace(clientSet, "phm")
	log.Println(podsList)
	podName := ""
	for _, item := range podsList {
		arr := strings.Split(item, "-deploy")
		log.Println(arr[0])
		if arr[0] == project {
			podName = item
			break
		}
	}
	if podName == "" {
		return "", errors.New("pod not found")
	} else {
		return podName, nil
	}
}

func GetContainerName(project string, container string) (string, error) {
	switch container {
	case "jupyter":
		return project + "-container", nil
	case "mlflow":
		return project + "-contianer-mlflow", nil
	default:
		return "", errors.New("No such container")
	}
}

func GetPhmFileList(c *gin.Context) {
	clientSet := utils.ClientSet
	var req GetFileListReq
	if c.ShouldBind(&req) == nil {
		log.Println(req)
	}
	podName, err := GetPodNameFromProject(req.ProjectName, clientSet)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	containerName, err := GetContainerName(req.ProjectName, req.ContainerName)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	result, err := resourceService.ExecCommandAndGetResult(clientSet, "phm", podName, containerName, getFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:ExecError:" + err.Error(),
		})
		return
	}
	var fileList []FileItem
	reg := regexp.MustCompile(`(.{10})[ ]{1,}([0-9])[ ]{1,}(\S*)[ ]{1,}(\S*)[ ]{1,}([0-9]*)[ ]{1,}([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}UTC)[ ]{1,}(.*)`)
	for i, line := range strings.Split(result, "\n") {
		if i == 0 {
			continue
		}
		// fmt.Println(line)
		arr := reg.FindStringSubmatch(line)
		if len(arr) != 8 {
			continue
		}
		size, _ := strconv.Atoi(arr[5])
		mtime, _ := time.Parse("2006-01-02T15:04:05UTC", arr[6])
		// fmt.Println(err)
		// fmt.Println(size, mtime, arr[6])
		fileList = append(fileList, FileItem{Name: arr[7], Size: size, ModifiedTime: &mtime})
		// fmt.Println(arr[4], arr[5], arr[6])

	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Data:       fileList,
	})
}

type GetPhmExecResultReq struct {
	ProjectName   string   `json:"project_name"`
	Command       []string `json:"command"`
	ContainerName string   `json:"container_name"`
}

func GetPhmExecResult(c *gin.Context) {
	clientSet := utils.ClientSet
	var req GetPhmExecResultReq
	if c.ShouldBind(&req) == nil {
		log.Println(req)
	}
	podName, err := GetPodNameFromProject(req.ProjectName, clientSet)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	containerName, err := GetContainerName(req.ProjectName, req.ContainerName)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	result, err := resourceService.ExecCommandAndGetResult(clientSet, "phm", podName, containerName, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:ExecError:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Data:       result,
	})
}

func GetPhmFile(c *gin.Context) {
	clientSet := utils.ClientSet
	var req GetFileReq
	if c.ShouldBind(&req) == nil {
		log.Println(req)
	}
	podName, err := GetPodNameFromProject(req.ProjectName, clientSet)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	containerName, err := GetContainerName(req.ProjectName, req.ContainerName)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	reader, writer := io.Pipe()
	resourceService.ExecCommandAndCopy(clientSet, "phm", podName, containerName, req.FileName, writer)
	// buf, err := io.ReadAll(reader)
	if err != nil {
		// fmt.Println(err)
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:CopyError:" + err.Error(),
		})
		return
	}
	filenameNew := fmt.Sprintf("%s/%s-%v.tar", req.ProjectName, req.FileName, time.Now().Unix())
	result, err := ossService.Upload(filenameNew, reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:UploadError:" + err.Error(),
		})
		return
	}
	// fmt.Println(base64.StdEncoding.EncodeToString(buf))
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Data: OssFileResp{
			URL:            result,
			OriginFileName: req.FileName,
		},
	})
}

type GetPodStatusReq struct {
	PodName string `json:"podname"`
}
type PodStatusResp struct {
	common.BaseResponse
	Data v1.ContainerStatus
}

func GetIaiPodStatus(c *gin.Context) {
	var req GetPodStatusReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	}
	clientSet := utils.ClientSet
	if pod, err := resourceService.GetPodInfo(clientSet, "iai", req.PodName); err != nil {
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	} else {
		for _, container := range pod.Status.ContainerStatuses {
			if container.Name == "runtime" {
				c.JSON(http.StatusOK, PodStatusResp{
					BaseResponse: common.BaseResponse{
						StatusCode: 0,
						StatusMsg:  "success",
					},
					Data: container,
				})
				break
			}
		}
	}
}

type GetPodNameReq struct {
	Taskid string `json:"taskid"`
	Podid  string `json:"podid"`
}
type PodNameResp struct {
	common.BaseResponse
	Data string
}

func GetIaiPodName(c *gin.Context) {
	var req GetPodNameReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	}
	clientSet := utils.ClientSet
	if pod, err := resourceService.GetPodBySelector(clientSet, "iai", req.Podid, req.Taskid); err != nil {
		c.JSON(http.StatusInternalServerError, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, PodNameResp{
			BaseResponse: common.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			Data: pod.Name,
		})
	}
}

type GetIaiPodHostNameReq struct {
	Podid string `json:"podid"`
}

func GetIaiPodHostName(c *gin.Context) {
	var req GetIaiPodHostNameReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	}
	var db = utils.DB
	var podInfo model.Pod
	res := db.Where("pod_id = ?", req.Podid).First(&podInfo)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + res.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Data:       podInfo.Node,
	})
}
