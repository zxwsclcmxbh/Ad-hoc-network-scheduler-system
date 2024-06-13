package taskService

import (
	"cloud/brainController/model"
	"cloud/brainController/utils"
	"errors"
	"time"
)

type TaskInfo struct {
	ID           int       `gorm:"column:id"`
	TaskName     string    `gorm:"column:task_name"`
	CreatedTime  time.Time `gorm:"column:created_time"`
	TaskDescribe string    `gorm:"column:task_describe"`
}

func GetPod(taskid string, podid string) (model.Pod, error) {
	var db = utils.DB
	var result = model.Pod{}
	res := db.Where(&model.Pod{TaskId: taskid, PodId: podid}).Find(&result)
	if res.RowsAffected != 1 {
		return result, errors.New("not found")
	}
	return result, res.Error
}

func GetSrcQueue(taskid string, podid string) ([]model.QueueRoute, error) {
	var db = utils.DB
	var result = []model.QueueRoute{}
	res := db.Where(&model.QueueRoute{TaskId: taskid, SrcPod: podid}).Find(&result)
	return result, res.Error
}

func GetDstQueue(taskid string, podid string) ([]model.QueueRoute, error) {
	var db = utils.DB
	var result = []model.QueueRoute{}
	res := db.Where(&model.QueueRoute{TaskId: taskid, DstPod: podid}).Find(&result)
	return result, res.Error
}

func GetTaskList(uid ...string) []model.Task {
	var db = utils.DB
	var tasks []model.Task
	if len(uid) == 0 {
		db.Select("created_at", "task_name", "task_describe", "task_status", "task_id", "intake_port").Find(&tasks)
	} else {
		db.Select("created_at", "task_name", "task_describe", "task_status", "task_id", "intake_port").Where("user_id=?", uid[0]).Find(&tasks)
	}
	return tasks
}

func GetTask(taskId string, uid string) (model.Task, error) {
	var db = utils.DB
	var task = model.Task{}
	result := db.Where(&model.Task{TaskId: taskId, UserId: uid}).Preload("Pods").Preload("Queues").Find(&task)
	if result.RowsAffected != 1 {
		return task, errors.New("not found")
	}
	return task, result.Error
}

func GetTasks(uid string) ([]AllTask, error) {
	var db = utils.DB
	var task = []model.Task{}
	var result = []AllTask{}
	err := db.Table("tasks").Where(&model.Task{UserId: uid}).Preload("Pods").Find(&task).Error
	for _, item := range task {
		t := []Pod{}
		for _, pod := range item.Pods {
			t = append(t, Pod{
				PodId: pod.PodId,
				Node:  pod.Node,
				Type:  pod.Type,
			})
		}
		result = append(result, AllTask{
			TaskName: item.TaskName,
			TaskId:   item.TaskId,
			Pods:     t,
		})
	}
	return result, err
}

type Pod struct {
	PodId string
	Node  string
	Type  string
}
type AllTask struct {
	TaskName string
	TaskId   string
	Pods     []Pod
}
