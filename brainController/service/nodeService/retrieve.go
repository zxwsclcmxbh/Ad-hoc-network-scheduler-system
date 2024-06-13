package nodeService

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetNodeInfo 获取节点信息
func GetNodeInfo(clientSet *kubernetes.Clientset, nodeName string) (map[string]interface{}, error) {
	nodeRel, err := clientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nodeInfo := map[string]interface{}{
		"createTime":    nodeRel.CreationTimestamp,
		"NowTime":       nodeRel.Status.Conditions[0].LastHeartbeatTime,
		"kernelVersion": nodeRel.Status.NodeInfo.KernelVersion,
		"systemOs":      nodeRel.Status.NodeInfo.OSImage,
		"cpu":           nodeRel.Status.Capacity.Cpu(),
		"docker:":       nodeRel.Status.NodeInfo.ContainerRuntimeVersion,
		"status":        nodeRel.Status.Conditions[len(nodeRel.Status.Conditions)-1].Type,
		"mem":           nodeRel.Status.Allocatable.Memory().String(),
	}
	return nodeInfo, nil
}

// GetNodeList 获取节点列表
func GetNodeList(clientSet *kubernetes.Clientset) ([]v1.Node, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

type NodeInfoNew struct {
	Metadata metav1.ObjectMeta    `json:"metadata"`
	Usage    *v1.ResourceList     `json:"usage"`
	Capacity v1.ResourceList      `json:"capacity"`
	Status   v1.NodeConditionType `json:"status"`
}

// GetNodeList 获取节点列表
func GetNodeListWithMetrics(clientSet *metricsv.Clientset, coreclient *kubernetes.Clientset) ([]NodeInfoNew, error) {
	nodesmetrics, err := clientSet.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes, err := coreclient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	// var resource = make(map[string]v1.ResourceList)
	var metrics = make(map[string]v1beta1.NodeMetrics)
	// var status = make(map[string]v1.NodeStatus)
	for _, n := range nodesmetrics.Items {
		// resource[n.Name] = n.Status.Capacity
		metrics[n.Name] = n
	}
	var result []NodeInfoNew
	for _, n := range nodes.Items {
		t := NodeInfoNew{}
		if metric, ok := metrics[n.Name]; ok {
			t = NodeInfoNew{
				Metadata: n.ObjectMeta,
				Usage:    &metric.Usage,
				Capacity: n.Status.Capacity,
				Status:   n.Status.Conditions[len(n.Status.Conditions)-1].Type,
			}
		} else {
			t = NodeInfoNew{
				Metadata: n.ObjectMeta,
				Usage:    nil,
				Capacity: n.Status.Capacity,
				Status:   n.Status.Conditions[len(n.Status.Conditions)-1].Type,
			}
		}
		result = append(result, t)
	}
	return result, nil
}
