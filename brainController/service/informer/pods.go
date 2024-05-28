package informer

import (
	"cloud/brainController/model"
	"cloud/brainController/service/cerebellumService"
	"cloud/brainController/service/taskService"
	"cloud/brainController/utils"
	"log"

	"time"

	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var CerebellumCh chan model.MessageWithNode

func InitInformer() {
	factory := informers.NewSharedInformerFactory(utils.ClientSet, time.Hour*24)
	podInformer := factory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{FilterFunc: filter, Handler: cache.ResourceEventHandlerFuncs{UpdateFunc: onUpdate}})
	stopchn := make(chan struct{})
	defer close(stopchn)
	factory.Start(stopchn)
	if !cache.WaitForCacheSync(stopchn, podInformer.Informer().HasSynced) {
		log.Println("failed to sync")
	}
	CerebellumCh = make(chan model.MessageWithNode)
	defer close(CerebellumCh)
	go cerebellumService.SendMessage(CerebellumCh)
	select {}
}

func filter(obj interface{}) bool {
	pod := obj.(*v1.Pod)
	if pod.Namespace == "iai" {
		return true
	} else {
		return false
	}
}

func onDelete(obj interface{}) {

	pod := obj.(*v1.Pod)
	// if pod.DeletionTimestamp != nil && pod.DeletionTimestamp.Time.Before(time.Now()) {
	// 	return
	// }
	podid, podidok := pod.Labels["app"]
	taskid, taskidok := pod.Labels["task"]
	if !podidok || !taskidok {
		return
	}
	if _, err := taskService.GetPod(taskid, podid); err != nil {
		return
	} else {
		taskService.UpdatePodNode(podid, "")
		log.Println("delete-update-pod", podid)
		queues, _ := taskService.GetSrcQueue(taskid, podid)
		for _, queue := range queues {
			if queue.IsRemote {
				CerebellumCh <- model.MessageWithNode{Node: queue.SrcNode, Msg: model.DeleteRouteReq{queue.QueueId}}
			}
			queue.SrcNode = ""
			queue.IsRemote = false
			log.Println("delete-update-src", queue)
			taskService.UpdateQueue(&queue)
		}
		queues, _ = taskService.GetDstQueue(taskid, podid)
		for _, queue := range queues {
			if queue.IsRemote {
				CerebellumCh <- model.MessageWithNode{Node: queue.SrcNode, Msg: model.DeleteRouteReq{queue.QueueId}}
			}
			queue.DstNode = ""
			queue.IsRemote = false
			log.Println("delete-update-dst", queue)
			taskService.UpdateQueue(&queue)

		}
	}
}

func onUpdate(oldObj, newObj interface{}) {
	newpod := newObj.(*v1.Pod)
	podid, podidok := newpod.Labels["app"]
	taskid, taskidok := newpod.Labels["task"]
	if newpod.Status.Phase != v1.PodRunning {
		return
	}
	if newpod.DeletionTimestamp != nil {
		if newpod.DeletionTimestamp.Time.After(time.Now()) {
			onDelete(oldObj)
			return
		} else {
			log.Println("update-notupdate", podid, "is deleting")
			return
		}
	}

	if !podidok || !taskidok {
		return
	}
	if res, err := taskService.GetPod(taskid, podid); err != nil {
		return
	} else {
		if res.Node != newpod.Spec.NodeName {
			log.Println("update-update-pod", podid, res.Node, newpod.Spec.NodeName)
			taskService.UpdatePodNode(podid, newpod.Spec.NodeName)

			queues, _ := taskService.GetSrcQueue(taskid, podid)
			for _, queue := range queues {
				queue.SrcNode = newpod.Spec.NodeName
				queue.IsRemote = queue.SrcNode != queue.DstNode
				log.Println("update-update-src", queue)
				taskService.UpdateQueue(&queue)
				if queue.IsRemote && queue.DstNode != "" {
					CerebellumCh <- model.MessageWithNode{Node: queue.SrcNode, Msg: model.AddRouteReq{
						TaskId: queue.TaskId,
						Dst:    queue.QueueId,
						Src:    queue.QueueId,
						Node:   queue.DstNode,
						Svc:    "cerebellum-svc-" + queue.DstNode,
					}}
				}
			}
			queues, _ = taskService.GetDstQueue(taskid, podid)
			for _, queue := range queues {
				queue.DstNode = newpod.Spec.NodeName
				queue.IsRemote = queue.SrcNode != queue.DstNode
				taskService.UpdateQueue(&queue)
				log.Println("update-update-dst", queue)
				if queue.IsRemote && queue.SrcNode != "" {
					CerebellumCh <- model.MessageWithNode{Node: queue.SrcNode, Msg: model.AddRouteReq{
						TaskId: queue.TaskId,
						Dst:    queue.QueueId,
						Src:    queue.QueueId,
						Node:   queue.DstNode,
						Svc:    "cerebellum-svc-" + queue.DstNode,
					}}
				}
			}
		} else {
			return
		}

	}
}
