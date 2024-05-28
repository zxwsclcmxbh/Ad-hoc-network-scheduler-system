package taskService

import (
	"cloud/brainController/model"
	"cloud/brainController/utils"
)

func DeleteTask(task model.Task) error {
	var db = utils.DB
	return db.Select("Queues", "Pods").Delete(&task).Error
}
