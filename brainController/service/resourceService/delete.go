package resourceService

import (
	"context"
	"errors"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeleteSvcAndDeploy 删除指定service和deployment
func DeleteSvcAndDeploy(clientSet *kubernetes.Clientset, namespace string, deploymentName string, serviceName string) (bool, error) {
	emptyDeleteOptions := metav1.DeleteOptions{}

	// 删除service
	deleteSvcErr := clientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, emptyDeleteOptions)

	// 删除deployment(deployment管理pod)
	deleteDeployErr := clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, emptyDeleteOptions)
	if deleteSvcErr == nil && deleteDeployErr == nil {
		return true, nil
	}
	return false, errors.New(deleteSvcErr.Error() + " and " + deleteDeployErr.Error())

}

// DeleteIaiPodsAndSvc 删除指定service和deployment
func DeleteIaiPodsAndSvc(clientSet *kubernetes.Clientset, namespace string, taskid string) (bool, error) {

	deployments, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "task=" + taskid})
	if err != nil {
		log.Println(err)
	}
	services, err := clientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "task=" + taskid})
	if err != nil {
		log.Println(err)
	}

	for _, item := range deployments.Items {
		log.Println("delete deployments", item.Name, clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), item.Name, metav1.DeleteOptions{}))
	}

	for _, item := range services.Items {
		log.Println("delete svc", item.Name, clientSet.CoreV1().Services(namespace).Delete(context.TODO(), item.Name, metav1.DeleteOptions{}))
	}

	return true, nil
}

func DeleteIaiDeployment(clientSet *kubernetes.Clientset, namespace string, deploymentName string) (bool, error) {
	emptyDeleteOptions := metav1.DeleteOptions{}

	if err := clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, emptyDeleteOptions); err != nil {
		log.Println(err)
	}
	return true, nil
}
