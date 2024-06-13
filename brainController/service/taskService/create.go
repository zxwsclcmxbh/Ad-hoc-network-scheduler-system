package taskService

import (
	"cloud/brainController/model"
	"cloud/brainController/utils"
)

func CreateTask(req model.TaskCreateRequest, taskId string, pods []*model.Pod, queues []*model.QueueRoute, uid string) error {
	var db = utils.DB
	task := model.Task{
		TaskCreateRequest: req,
		TaskStatus:        "creating",
		TaskId:            taskId,
		Pods:              pods,
		Queues:            queues,
		UserId:            uid,
	}
	err := db.Create(&task).Error
	if err != nil {
		return err
	}
	db.Save(&task)
	return nil
}

func CreateTaskOld(req model.TaskCreateRequestOld, taskId string, pods []*model.Pod, queues []*model.QueueRoute) error {
	var db = utils.DB
	task := model.TaskOld{
		TaskCreateRequestOld: req,
		TaskStatus:           "creating",
		TaskId:               taskId,
		Pods:                 pods,
		Queues:               queues,
	}

	err := db.Table("task").Create(&task).Error
	if err != nil {
		return err
	}
	db.Save(&task)
	return nil
}
