package taskService

import (
	"cloud/brainController/common"
	"cloud/brainController/model"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func genEnv(key string) string {
	return fmt.Sprintf("DYNACONF_%s", key)
}
func GenUUID() string {
	t := uuid.NewString()
	return strings.ReplaceAll(t, "-", "")
}
func GenPods(task []model.TaskDefinition) (taskId string, pods []*model.Pod, routes []*model.QueueRoute) {
	taskId = GenUUID()
	var tempQueue *model.QueueRoute
	if len(task) == 0 {
		return
	}
	if len(task) == 1 {
		podId := GenUUID()
		env := make(map[string]string)
		env[genEnv("TYPE")] = task[0].Image
		env[genEnv("INPUT")] = "input"
		env[genEnv("OUPUT")] = "output"
		env[genEnv("BASE_PATH")] = "/data/"
		env[genEnv("TASK_ID")] = taskId
		env[genEnv("PORT")] = "8080"
		env[genEnv("REDIS_HOST")] = "redis-svc" + "-" + task[0].Node
		env[genEnv("CUSTOM")] = task[0].CustomConf
		env[genEnv("NODE")] = task[0].Node
		p := &model.Pod{Image: task[0].Image, PodId: podId, Node: task[0].Node, Environment: env, TaskId: taskId, Type: task[0].Type}
		pods = append(pods, p)
		return
	}
	for index, item := range task {
		podid := GenUUID()
		env := make(map[string]string)
		env[genEnv("TYPE")] = item.Type
		q := &model.QueueRoute{
			SrcPod: podid, QueueId: GenUUID(), IsRemote: false, SrcNode: item.Node, TaskId: taskId,
		}
		p := &model.Pod{Image: item.Image, PodId: podid, Node: item.Node, Environment: env, TaskId: taskId, Type: item.Type}
		if index == 0 {
			env[genEnv("INPUT")] = "input"
			env[genEnv("OUTPUT")] = q.QueueId
			routes = append(routes, q)
		} else if index == len(task)-1 {
			tempQueue.DstPod = p.PodId
			env[genEnv("INPUT")] = tempQueue.QueueId
			tempQueue.DstNode = item.Node
			if tempQueue.SrcNode != item.Node {
				tempQueue.IsRemote = true
			}
			env[genEnv("OUTPUT")] = "output"
		} else {
			tempQueue.DstPod = p.PodId
			env[genEnv("INPUT")] = tempQueue.QueueId
			env[genEnv("OUTPUT")] = q.QueueId
			tempQueue.DstNode = item.Node
			if tempQueue.SrcNode != item.Node {
				tempQueue.IsRemote = true
			}
			routes = append(routes, q)
		}
		env[genEnv("REDIS_HOST")] = "redis-svc" + "-" + item.Node
		env[genEnv("NODE")] = item.Node
		env[genEnv("BASE_PATH")] = "/data/"
		env[genEnv("TASK_ID")] = taskId
		env[genEnv("PORT")] = "8080"
		env[genEnv("CUSTOM")] = item.CustomConf
		pods = append(pods, p)
		tempQueue = q
	}
	return
}

//	func main() {
//		tasks := []models.TaskDefination{{Image: "rolling-window:v1", Node: "node01", CustomConf: "{\"window_length\":128,\"slide_step\":1}"}, {Image: "pre-process:v1", Node: "node01", CustomConf: "{}"}, {Image: "feature-extraction:v1", Node: "node01", CustomConf: "{}"}, {Image: "bfd:v1", Node: "node01", CustomConf: "{}"}}
//		_, pods, routes := GenPods(tasks)
//		for _, item := range pods {
//			fmt.Println(item)
//		}
//		for _, item := range routes {
//			fmt.Println(item)
//		}
//	}
func getStartLinks(nodeid string, links []model.Link) []string {
	result := []string{}
	for _, link := range links {
		if link.StartId == nodeid {
			result = append(result, link.Id)
		}
	}
	return result
}
func getEndLinks(nodeid string, links []model.Link) []string {
	result := []string{}
	for _, link := range links {
		if link.EndId == nodeid {
			result = append(result, link.Id)
		}
	}
	return result
}

func getPod(podid string, pods []*model.Pod) *model.Pod {
	for _, pod := range pods {
		if pod.PodId == podid {
			return pod
		}
	}
	return nil
}

func formatListEnv(list []string) string {
	j, _ := json.Marshal(list)
	return "@json " + string(j)
}

func ComposeNew(taskdefination model.TaskDefination) (taskId string, pods []*model.Pod, routes []*model.QueueRoute, intakepodid []string) {
	taskId = GenUUID()
	for _, node := range taskdefination.Graph.NodeList {
		nodedata := taskdefination.Data.(map[string]interface{})[node.Id].(map[string]interface{})
		nodeAffinity := nodedata["nodeAffinity"].(string)
		var image string
		if index := strings.Index(node.Meta.Id, "models"); index != -1 {
			image = nodedata["models"].(string)
		} else {
			image = common.CommonImgae[node.Meta.Id]
		}
		intake_mode := ""
		if index := strings.Index(node.Meta.Id, "intake"); index != -1 {
			switch node.Meta.Id {
			case "intake-kafka":
				intake_mode = "kafka"
			case "intake-http-form":
				intake_mode = "http-file"
			case "intake-http-json":
				intake_mode = "http-json"
			}
			intakepodid = append(intakepodid, node.Id)
		}
		nodedata["mode"] = intake_mode
		env := make(map[string]string)
		custom_conf, _ := json.Marshal(nodedata)
		env[genEnv("TYPE")] = node.Meta.Id
		env[genEnv("INPUT")] = formatListEnv(getEndLinks(node.Id, taskdefination.Graph.LinkList))
		env[genEnv("OUTPUT")] = formatListEnv(getStartLinks(node.Id, taskdefination.Graph.LinkList))
		env[genEnv("BASE_PATH")] = "/data/"
		env[genEnv("POD")] = node.Id
		env[genEnv("TASK_ID")] = taskId
		env[genEnv("PORT")] = "8080"
		// env[genEnv("REDIS_HOST")] = "redis-svc" + "-" + nodename
		env[genEnv("CUSTOM")] = string(custom_conf)
		// env[genEnv("NODE")] = nodename
		pod := &model.Pod{
			PodId:        node.Id,
			TaskId:       taskId,
			Node:         "",
			Type:         node.Meta.Id,
			Image:        image,
			Environment:  env,
			NodeAffinity: nodeAffinity,
		}
		pods = append(pods, pod)
	}
	for _, link := range taskdefination.Graph.LinkList {
		podstart := getPod(link.StartId, pods)
		podend := getPod(link.EndId, pods)
		q := model.QueueRoute{
			QueueId: link.Id,
			TaskId:  taskId,
			SrcPod:  podstart.PodId,
			DstPod:  podend.PodId,
		}
		routes = append(routes, &q)
	}
	return
}

func RecomposePod(env map[string]string, node string) {
	env[genEnv("REDIS_HOST")] = "redis-svc" + "-" + node
	env[genEnv("NODE")] = node
}
