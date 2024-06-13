package controller

import (
	"cloud/brainController/common"
	"cloud/brainController/model"
	"cloud/brainController/service/cerebellumService"
	"cloud/brainController/service/resourceService"
	"cloud/brainController/service/taskService"
	"cloud/brainController/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetTaskListResponse struct {
	common.BaseResponse
	Data []model.Task `json:"data"`
}

type CreateTaskResponse struct {
	common.BaseResponse
	Data struct {
		TaskId string `json:"task_id"`
	}
}

func CreateTask(c *gin.Context) {
	var req model.TaskCreateRequestOld
	err := c.ShouldBind(&req)
	log.Println(req)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	taskId, pods, queues := taskService.GenPods(req.TaskDefinition)
	log.Println(pods, queues)
	for _, q := range queues {
		if q.IsRemote {
			msg := model.AddRouteReq{
				TaskId: q.TaskId,
				Dst:    q.QueueId,
				Src:    q.QueueId,
				Node:   q.DstNode,
				Svc:    "cerebellum-svc-" + q.DstNode,
			}
			log.Println(cerebellumService.SendRoute(msg, q.SrcNode))
		}
	}
	err = taskService.CreateTaskOld(req, taskId, pods, queues)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	for _, pod := range pods {
		// TODO: rollback?
		NodeAffinity := ""
		fmt.Print(resourceService.CreateIaiDeployment(utils.ClientSet, "iai", pod.Type+"-"+pod.PodId, 1, pod.PodId, pod.Node, "runtime", pod.Image, 8080, taskId, utils.Transform(pod.Environment), NodeAffinity))
	}
	// TODO: rollback?
	// fmt.Println(resourceService.CreateIaiService(utils.ClientSet, "iai", taskId, 8080, pods[0].PodId))
	// go taskService.CheckStatus(pods, taskId)
	c.JSON(http.StatusOK, CreateTaskResponse{
		Data: struct {
			TaskId string `json:"task_id"`
		}{TaskId: taskId},
		BaseResponse: common.BaseResponse{
			StatusCode: 200,
			StatusMsg:  "success:CreateTask",
		},
	})
}

type PodImagesUpdateRequest struct {
	TaskID     string `json:"task_id"`
	PodID      string `json:"pod_id"`
	NewImage   string `json:"new_image"`
	Definition string `json:"definition"`
}

func PodImagesUpdate(c *gin.Context) {
	var req PodImagesUpdateRequest
	err := c.ShouldBind(&req)
	fmt.Printf("req:%#v\n", req)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	task, err := taskService.GetTask(req.TaskID, c.GetString("uid"))
	fmt.Printf("task:%#v\n", task)
	if err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	}
	if task.TaskStatus != "UP" {
		c.JSON(http.StatusNotAcceptable, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + "task status is not up",
		})
		return
	}
	var temppod *model.Pod
	for _, pod := range task.Pods {
		if pod.PodId == req.PodID {
			temppod = pod
			break
		}
	}
	if temppod == nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed: can't find pod",
		})
		return
	}
	log.Println(taskService.UpdateTaskStatusAndDefinition(req.TaskID, "updating", req.Definition))
	_, err = resourceService.UpdateIaiDeploymentImage(utils.ClientSet, "iai", temppod.Type+"-"+temppod.PodId+"-"+temppod.TaskId, req.NewImage)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	err = taskService.UpdatePodImage(req.PodID, req.NewImage)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	go taskService.CheckMigrationStatus(req.PodID, req.TaskID)
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

type PodMigrationRequest struct {
	TaskID     string `json:"task_id"`
	PodID      string `json:"pod_id"`
	TargetNode string `json:"target_node"`
	Definition string `json:"definition"`
}

func PodMigration(c *gin.Context) {
	var req PodMigrationRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	if task, err := taskService.GetTask(req.TaskID, c.GetString("uid")); err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	} else {
		if task.TaskStatus != "UP" {
			c.JSON(http.StatusNotAcceptable, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "failed:" + "task status is not up",
			})
			return
		}
		var temppod *model.Pod
		for _, pod := range task.Pods {
			if pod.PodId == req.PodID {
				temppod = pod
				break
			}
		}
		if temppod == nil {
			c.JSON(http.StatusNotFound, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "failed:" + err.Error(),
			})
			return
		}
		// } else {
		// 	var start, end *models.QueueRoute
		// 	for _, queue := range task.Queues {
		// 		if queue.SrcPod == req.PodID {
		// 			queue := queue
		// 			start = queue
		// 			break
		// 		}
		// 	}
		// 	for _, queue := range task.Queues {
		// 		if queue.DstPod == req.PodID {
		// 			queue := queue
		// 			end = queue
		// 			break
		// 		}
		// 	}
		// 	log.Println("start", start)
		// 	log.Println("end", end)
		// 	if start != nil {

		// 		if start.IsRemote {
		// 			cerebellumService.DeleteRoute(models.DeleteRouteReq{RouteId: start.QueueId}, start.SrcNode)
		// 		}
		// 		start.SrcNode = req.TargetNode
		// 		start.IsRemote = start.SrcNode != start.DstNode
		// 		if start.IsRemote {
		// 			r := models.AddRouteReq{TaskId: req.TaskID, Dst: start.QueueId, Src: start.QueueId, Node: start.DstNode, Svc: "cerebellum-svc-" + start.DstNode}
		// 			log.Println("start", r)
		// 			cerebellumService.SendRoute(r, start.SrcNode)
		// 		}
		// 		log.Println(taskService.UpdateQueue(start))
		// 	}
		// 	if end != nil {
		// 		end.DstNode = req.TargetNode
		// 		if end.IsRemote {
		// 			cerebellumService.DeleteRoute(models.DeleteRouteReq{RouteId: end.QueueId}, end.SrcNode)
		// 		}
		// 		end.IsRemote = end.SrcNode != end.DstNode
		// 		if end.IsRemote {
		// 			r := models.AddRouteReq{TaskId: req.TaskID, Dst: end.QueueId, Src: end.QueueId, Node: end.DstNode, Svc: "cerebellum-svc-" + end.DstNode}
		// 			log.Println("end", r)
		// 			cerebellumService.SendRoute(r, end.SrcNode)
		// 		}
		// 		log.Println(taskService.UpdateQueue(end))
		// 	}

		log.Println(taskService.UpdateTaskStatusAndDefinition(req.TaskID, "migrating", req.Definition))
		// log.Println(taskService.UpdatePodNode(req.PodID, req.TargetNode))
		resourceService.DeleteIaiDeployment(utils.ClientSet, "iai", temppod.Type+"-"+temppod.PodId+"-"+temppod.TaskId)
		// log.Println(temppod.Environment)
		// taskService.RecomposePod(temppod.Environment, req.TargetNode)
		// log.Println(temppod.Environment)
		log.Println(resourceService.CreateIaiDeployment(utils.ClientSet, "iai", temppod.Type+"-"+temppod.PodId, 1, temppod.PodId, req.TargetNode, "runtime", temppod.Image, 8080, req.TaskID, utils.Transform(temppod.Environment), temppod.NodeAffinity))
		go taskService.CheckMigrationStatus(req.PodID, req.TaskID)
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	}
}

type nodeAffinityMigrationRequest struct {
	TaskID          string `json:"task_id"`
	PodID           string `json:"pod_id"`
	NewNodeAffinity string `json:"new_nodeAffinity"`
	Definition      string `json:"definition"`
}

func NodeAffinityMigration(c *gin.Context) {
	var req nodeAffinityMigrationRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	if task, err := taskService.GetTask(req.TaskID, c.GetString("uid")); err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	} else {
		if task.TaskStatus != "UP" {
			c.JSON(http.StatusNotAcceptable, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "failed:" + "task status is not up",
			})
			return
		}
		var temppod *model.Pod
		for _, pod := range task.Pods {
			if pod.PodId == req.PodID {
				temppod = pod
				break
			}
		}
		if temppod == nil {
			c.JSON(http.StatusNotFound, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "failed:" + err.Error(),
			})
			return
		}
		log.Println(taskService.UpdateTaskStatusAndDefinition(req.TaskID, "migrating", req.Definition))
		_, err := resourceService.DeleteIaiDeployment(utils.ClientSet, "iai", temppod.Type+"-"+temppod.PodId+"-"+temppod.TaskId)
		if err != nil {
			c.JSON(http.StatusNotFound, common.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "failed:" + err.Error(),
			})
			return
		}
		log.Println(resourceService.CreateIaiDeployment(utils.ClientSet, "iai", temppod.Type+"-"+temppod.PodId, 1, temppod.PodId, "", "runtime", temppod.Image, 8080, req.TaskID, utils.Transform(temppod.Environment), req.NewNodeAffinity))
		go taskService.CheckMigrationStatus(req.PodID, req.TaskID)
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	}
}

func GetTaskList(c *gin.Context) {
	TaskListInfo := taskService.GetTaskList(c.GetString("uid"))
	// fmt.Println(TaskListInfo)
	c.JSON(http.StatusOK, GetTaskListResponse{
		Data: TaskListInfo,
		BaseResponse: common.BaseResponse{
			StatusCode: 200,
			StatusMsg:  "success:GetTaskList",
		},
	})
}

func GetTask(c *gin.Context) {
	taskid := c.Query("taskid")

	if task, err := taskService.GetTask(taskid, c.GetString("uid")); err != nil {
		c.JSON(http.StatusNotFound, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			Data:       task,
		})
	}
}

type DeleteTaskReq struct {
	TaskID string `json:"task_id"`
}

func DeleteTask(c *gin.Context) {
	var r DeleteTaskReq
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	if task, err := taskService.GetTask(r.TaskID, c.GetString("uid")); err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	} else {
		taskService.DeleteTask(task)

		for _, i := range task.Queues {
			if i.IsRemote {
				log.Println(cerebellumService.DeleteRoute(model.DeleteRouteReq{i.QueueId}, i.SrcNode))
			}
		}
		resourceService.DeleteIaiPodsAndSvc(utils.ClientSet, "iai", task.TaskId)

	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 200,
		StatusMsg:  "success",
	})
}

func CreateTaskNew(c *gin.Context) {
	var req model.TaskCreateRequest
	err := c.ShouldBind(&req)
	log.Println(req)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	var task model.TaskDefination
	json.Unmarshal([]byte(req.TaskDefinition), &task)
	taskId, pods, queues, intakepodid := taskService.ComposeNew(task)
	// for _, pod := range pods {
	// 	log.Println(pod)
	// }
	// for _, queue := range queues {
	// 	log.Println(queue)
	// }
	// for _, q := range queues {
	// 	if q.IsRemote {
	// 		msg := models.AddRouteReq{
	// 			TaskId: q.TaskId,
	// 			Dst:    q.QueueId,
	// 			Src:    q.QueueId,
	// 			Node:   q.DstNode,
	// 			Svc:    "cerebellum-svc-" + q.DstNode,
	// 		}
	// 		log.Println(cerebellumService.SendRoute(msg, q.SrcNode))
	// 	}
	// }
	err = taskService.CreateTask(req, taskId, pods, queues, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	for _, pod := range pods {
		// TODO: rollback?
		fmt.Println("nodeAffinity")
		fmt.Printf("%#v\n\n", pod.NodeAffinity)
		fmt.Print(resourceService.CreateIaiDeployment(utils.ClientSet, "iai", pod.Type+"-"+pod.PodId, 1, pod.PodId, pod.Node, "runtime", pod.Image, 8080, taskId, utils.Transform(pod.Environment), pod.NodeAffinity))
	}
	// TODO: rollback?
	var intakeids []string
	for index, podid := range intakepodid {
		fmt.Println(resourceService.CreateIaiService(utils.ClientSet, "iai", taskId, strconv.Itoa(index), 8080, podid))
		intakeids = append(intakeids, taskId+"-"+strconv.Itoa(index))
	}
	go taskService.CheckStatus(pods, taskId, intakeids)
	c.JSON(http.StatusOK, CreateTaskResponse{
		Data: struct {
			TaskId string `json:"task_id"`
		}{TaskId: taskId},
		BaseResponse: common.BaseResponse{
			StatusCode: 200,
			StatusMsg:  "success:CreateTask",
		},
	})

}

func GetAllTasks(c *gin.Context) {
	tasks, err := taskService.GetTasks(c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: 1,
			StatusMsg:  "failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Data:       tasks,
	})
}
