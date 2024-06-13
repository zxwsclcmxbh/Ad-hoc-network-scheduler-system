package taskService

import (
	"cloud/brainController/model"
	"cloud/brainController/service/resourceService"
	"cloud/brainController/utils"
	"fmt"
	"log"
	"sync"
)

func UpdateTask(taskId string, port []int, status string) error {
	var db = utils.DB
	var task model.Task
	db.Where(&model.Task{TaskId: taskId}).Find(&task)
	task.IntakePort = port
	task.TaskStatus = status
	return db.Save(&task).Error
}

func CheckStatus(pods []*model.Pod, taskId string, svcids []string) {
	wg := new(sync.WaitGroup)
	for _, pod := range pods {
		wg.Add(1)
		go func(podId string) {
			resourceService.WaitDeploymentReady(utils.ClientSet, "iai", podId)
			wg.Done()
		}(pod.PodId)
	}
	wg.Wait()
	ports := []int{}
	for _, svcid := range svcids {
		svc, _ := resourceService.GetSvcInfo(utils.ClientSet, "iai", "intake-"+svcid)
		ports = append(ports, int(svc.Spec.Ports[0].NodePort))
	}
	log.Printf("%s deploy finished\n", taskId)
	UpdateTask(taskId, ports, "UP")
}

func UpdateTaskStatusAndDefinition(taskId string, status string, definition string) error {
	var db = utils.DB
	var task model.Task
	db.Where(&model.Task{TaskId: taskId}).Find(&task)
	task.TaskStatus = status
	task.TaskDefinition = definition
	return db.Save(&task).Error
}

func UpdateTaskStatus(taskId string, status string) error {
	var db = utils.DB
	var task model.Task
	db.Where(&model.Task{TaskId: taskId}).Find(&task)
	task.TaskStatus = status
	return db.Save(&task).Error
}

func UpdatePodNode(podid string, node string) error {
	var db = utils.DB
	var pod model.Pod
	db.Where(&model.Pod{PodId: podid}).Find(&pod)
	pod.Node = node
	return db.Save(&pod).Error
}

func UpdatePodImage(podid string, images string) error {
	var db = utils.DB
	var pod model.Pod
	db.Where(&model.Pod{PodId: podid}).Find(&pod)
	fmt.Printf("pod-xbh:%v", pod)
	pod.Image = images
	fmt.Printf("pod-xbh:%v", pod)
	return db.Save(&pod).Error
}

func UpdateQueue(queue *model.QueueRoute) error {
	var db = utils.DB
	return db.Save(queue).Error
}

func CheckMigrationStatus(podId string, taskId string) {
	resourceService.WaitDeploymentReady(utils.ClientSet, "iai", podId)
	log.Printf("%s deploy finished\n", taskId)
	UpdateTaskStatus(taskId, "UP")
}
