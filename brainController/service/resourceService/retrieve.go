package resourceService

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// WatchContainerLogWithPodName
// 参数:
// Container:容器名称
// Follow:跟踪Pod的日志流，默认为false（关闭）对应kubectl logs命令中的 -f 参数
// TailLines:如果设置，则显示从日志末尾开始的行数。如果未指定，则从容器的创建开始或从秒开始或从时间开始显示日志
// Previous:返回以前终止的容器日志。默认为false（关闭）
func WatchContainerLogWithPodName(clientSet *kubernetes.Clientset, namespace string, podName string) (io.ReadCloser, error) {
	logOpt := &v1.PodLogOptions{
		Follow:   false,
		Previous: false,
	}
	req := clientSet.CoreV1().Pods(namespace).GetLogs(podName, logOpt)
	return req.Stream(context.TODO())
}

func GetNamespaces(clientSet *kubernetes.Clientset) (namespaceList []string) {
	list, _ := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	for _, item := range list.Items {
		namespaceList = append(namespaceList, item.Name)
	}
	return
}

func GetPodsListInNamespace(clientSet *kubernetes.Clientset, namespace string) []string {
	list, _ := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	var podsList []string
	for _, item := range list.Items {
		podsList = append(podsList, item.Name)
	}
	return podsList
}

func GetPodInfo(clientSet *kubernetes.Clientset, namespace string, podName string) (*v1.Pod, error) {
	podInfo, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	//return podInfo.Status.PodIP, nil
	//podInfo.Spec.Hostname
	return podInfo, nil
}

func GetPodBySelector(clientSet *kubernetes.Clientset, namespace string, podid string, taskid string) (*v1.Pod, error) {
	label := map[string]string{
		"task": taskid,
		"app":  podid,
	}
	podInfo, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: label})})
	if err != nil {
		return nil, err
	}
	return &podInfo.Items[0], nil
}

func GetSvcListInNamespace(clientSet *kubernetes.Clientset, namespace string) []string {
	list, _ := clientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	var svcList []string
	for _, item := range list.Items {
		svcList = append(svcList, item.Name)
	}
	return svcList
}

func GetSvcInfo(clientSet *kubernetes.Clientset, namespace string, svcName string) (*v1.Service, error) {
	svcInfo, err := clientSet.CoreV1().Services(namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return svcInfo, nil
}

func CreateDeploymentWatcher(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podId string) (watch.Interface, error) {
	opts := metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: "app=" + podId,
	}
	log.Println("watcher create for" + podId)
	return clientSet.AppsV1().Deployments(namespace).Watch(ctx, opts)
}

func WaitDeploymentReady(clientSet *kubernetes.Clientset, namespace string, podId string) error {
	ctx, _ := context.WithTimeout(context.TODO(), 500*time.Second)
	watcher, err := CreateDeploymentWatcher(ctx, clientSet, namespace, podId)
	defer watcher.Stop()
	if err != nil {
		return errors.New("unable to watch")
	} else {
		for {
			select {
			case event := <-watcher.ResultChan():
				// log.Println(event.Type)
				deployment := event.Object.(*appsv1.Deployment)
				// log.Println(deployment, deployment.Status)
				if deployment.Status.ReadyReplicas == 1 {
					return nil
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
