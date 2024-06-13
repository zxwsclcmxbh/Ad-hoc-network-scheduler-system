package resourceService

import (
	"cloud/brainController/common"
	"cloud/brainController/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/klauspost/cpuid/v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

// CreateService 新建service
func CreateService(clientSet *kubernetes.Clientset, namespace string, serviceName string, imagePort int32, selector string) (bool, error) {
	// 得到service的客户端
	serviceClient := clientSet.CoreV1().Services(namespace)

	// 实例化一个数据结构
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{{
				Name: "http",
				Port: imagePort,
			},
			},
			Selector: map[string]string{
				"app": selector,
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	log.Printf("Create services %s \n", result.GetName())
	return true, nil
}

func CreatePHMService(clientSet *kubernetes.Clientset, namespace string, serviceName string, imagePort int32, selector string) (bool, error) {
	// 得到service的客户端
	serviceClient := clientSet.CoreV1().Services(namespace)

	// 实例化一个数据结构
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{{
				Name: "http",
				Port: imagePort,
			}, {
				Name: "mlflow",
				Port: 8889,
			},
			},
			Selector: map[string]string{
				"app": selector,
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}

	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	log.Printf("Create services %s \n", result.GetName())
	return true, nil
}

// CreateIaiService 新建service
func CreateIaiService(clientSet *kubernetes.Clientset, namespace string, taskid string, serviceName string, imagePort int32, selector string) (bool, error) {
	// 得到service的客户端
	serviceClient := clientSet.CoreV1().Services(namespace)

	// 实例化一个数据结构
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "intake-" + taskid + "-" + serviceName,
			Labels: map[string]string{
				"task": taskid,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{{
				Name: "http",
				Port: imagePort,
			},
			},
			Selector: map[string]string{
				"app": selector,
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	log.Printf("Create services %s \n", result.GetName())
	return true, nil
}

// CreateDeployment 新建deployment
func CreateDeployment(
	clientSet *kubernetes.Clientset,
	namespace string,
	deploymentName string,
	replicas int32,
	matchLabels string,
	nodeName string,
	containerName string,
	imageName string,
	imagePort int32,
	isGPU bool) (bool, error) {
	// 得到deployment的客户端
	deploymentClient := clientSet.AppsV1().Deployments(namespace)
	//实力化container
	var containers []apiv1.Container
	var spec apiv1.PodSpec
	log.Println(isGPU)
	if isGPU {
		containers = []apiv1.Container{
			{
				Name:  containerName,
				Image: imageName,
				Resources: apiv1.ResourceRequirements{
					Limits: apiv1.ResourceList{
						//"nvidia.com/gpu": resource.MustParse("1"),
						"aliyun.com/gpu-mem": resource.MustParse("1"),
					},
				},
				ImagePullPolicy: "IfNotPresent",
				Ports: []apiv1.ContainerPort{
					{
						Name:          "http",
						Protocol:      apiv1.ProtocolSCTP,
						ContainerPort: imagePort,
					},
				},
			},
		}
		spec = apiv1.PodSpec{
			Containers: containers,
		}
	} else {
		containers = []apiv1.Container{
			{
				Name:            containerName,
				Image:           imageName,
				ImagePullPolicy: "IfNotPresent",
				Ports: []apiv1.ContainerPort{
					{
						Name:          "http",
						Protocol:      apiv1.ProtocolSCTP,
						ContainerPort: imagePort,
					},
				},
			},
		}
		spec = apiv1.PodSpec{
			NodeName:   nodeName,
			Containers: containers,
		}
	}
	log.Println("containers:")
	log.Println(containers)

	// 实例化一个数据结构
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: map[string]string{
				"app": matchLabels,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": matchLabels, // 用于select pod
				},
			},

			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": matchLabels,
					},
				},
				Spec: spec,
			},
		},
	}

	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return false, err
	}

	log.Printf("Create deployment %s \n", result.GetName())
	return true, nil
}

func CreatePHMDeployment(
	clientSet *kubernetes.Clientset,
	namespace string,
	deploymentName string,
	replicas int32,
	matchLabels string,
	nodeName string,
	containerName string,
	imageName string,
	imagePort int32,
	rid string,
	password string,
	hex string,
	isGPU bool) (bool, error) {
	// 得到deployment的客户端

	deploymentClient := clientSet.AppsV1().Deployments(namespace)
	//实力化container
	var containers []apiv1.Container
	var spec apiv1.PodSpec
	log.Println(isGPU)
	baseContainer := apiv1.Container{
		Name:            containerName + "-mlflow",
		Image:           common.CommonImgae["mlflow"],
		ImagePullPolicy: "IfNotPresent",
		Ports: []apiv1.ContainerPort{
			{
				Name:          "mlflow",
				ContainerPort: 8889,
				Protocol:      apiv1.ProtocolSCTP,
			},
		},
		Command: []string{"mlflow", "server", "--host", "0.0.0.0", "--port", "8889", "--backend-store-uri", "sqlite:///mlflow.db", "--serve-artifacts", "--artifacts-destination", "s3://mlflow/" + rid + "/"},
		Env: []apiv1.EnvVar{
			{
				Name:  "AWS_ACCESS_KEY_ID",
				Value: "v18xIvb5d2KrXFoj",
			}, {
				Name:  "AWS_SECRET_ACCESS_KEY",
				Value: "8Ouj2hTFFtC7wN0oLukqEWnpQpvSH4wa",
			}, {
				Name:  "MLFLOW_S3_ENDPOINT_URL",
				Value: "http://" + utils.Config.Minio.URL,
			},
		},
	}
	if isGPU {
		containers = []apiv1.Container{
			{
				Name:  containerName,
				Image: imageName,
				Resources: apiv1.ResourceRequirements{
					Limits: apiv1.ResourceList{
						//"nvidia.com/gpu": resource.MustParse("1"),
						"aliyun.com/gpu-mem": resource.MustParse("1"),
					},
				},
				Command:         []string{"start-notebook.sh", "--NotebookApp.base_url=/" + rid + "/jupyter" + "/", "--NotebookApp.password=" + hex},
				ImagePullPolicy: "IfNotPresent",
				Ports: []apiv1.ContainerPort{
					{
						Name:          "http",
						Protocol:      apiv1.ProtocolSCTP,
						ContainerPort: imagePort,
					},
				},
			},
		}
		containers = append(containers, baseContainer)
		spec = apiv1.PodSpec{
			Containers: containers,
		}
	} else {
		containers = []apiv1.Container{
			{
				Name:            containerName,
				Image:           imageName,
				ImagePullPolicy: "IfNotPresent",
				Command:         []string{"start-notebook.sh", "--NotebookApp.base_url=/" + rid + "/jupyter" + "/", "--NotebookApp.password=" + hex},
				Ports: []apiv1.ContainerPort{
					{
						Name:          "http",
						Protocol:      apiv1.ProtocolSCTP,
						ContainerPort: imagePort,
					},
				},
			},
		}
		containers = append(containers, baseContainer)
		spec = apiv1.PodSpec{
			NodeName:   nodeName,
			Containers: containers,
		}
	}
	log.Println("containers:")
	log.Println(containers)

	// 实例化一个数据结构
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: map[string]string{
				"app": matchLabels,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": matchLabels, // 用于select pod
				},
			},

			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": matchLabels,
					},
				},
				//Spec: apiv1.PodSpec{
				//	NodeName:   nodeName,
				//	Containers: containers,
				//},
				Spec: spec,
			},
		},
	}

	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return false, err
	}

	log.Printf("Create deployment %s \n", result.GetName())
	return true, nil
}

func CreateIaiDeployment(
	clientSet *kubernetes.Clientset,
	namespace string,
	deploymentName string,
	replicas int32,
	matchLabels string,
	nodeName string,
	containerName string,
	imageName string,
	imagePort int32,
	taskId string,
	env []apiv1.EnvVar,
	nodeAffinity string) (bool, error) {
	// 得到deployment的客户端
	deploymentClient := clientSet.AppsV1().Deployments(namespace)
	//反序列化nodeAffinity
	nodeAffinityByte := []byte(nodeAffinity)
	var NodeAffinityStruct apiv1.NodeAffinity
	err := json.Unmarshal(nodeAffinityByte, &NodeAffinityStruct)
	if err != nil {
		return false, err
	}
	fmt.Printf("%#v\n\n", NodeAffinityStruct)
	// 实例化一个数据结构
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName + "-" + taskId,
			Labels: map[string]string{
				"task": taskId,
				"app":  matchLabels,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{Type: appsv1.RecreateDeploymentStrategyType},
			Replicas: pointer.Int32(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": matchLabels, // 用于select pod
				},
			},

			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  matchLabels,
						"task": taskId,
					},
				},
				Spec: apiv1.PodSpec{
					Affinity: &apiv1.Affinity{
						NodeAffinity: &NodeAffinityStruct,
					},
					Containers: []apiv1.Container{
						{
							Name:            containerName,
							Image:           imageName,
							ImagePullPolicy: "IfNotPresent",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolSCTP,
									ContainerPort: imagePort,
								},
							},
							Env: env,
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "data",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "data",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/data/" + taskId + "/",
									Type: (*apiv1.HostPathType)(pointer.String(string(apiv1.HostPathDirectoryOrCreate))),
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return false, err
	}

	log.Printf("Create deployment %s \n", result.GetName())
	return true, nil
}

func UpdateIaiDeploymentImage(clientSet *kubernetes.Clientset, namespace string, deploymentName string, newImageName string) (bool, error) {
	//获取deployment
	// deployName=pod.type-pod.id-task.id
	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	fmt.Printf("deployment:%#v\n", deployment)
	if err != nil {
		return false, err
	}
	//更改镜像
	deployment.Spec.Template.Spec.Containers[0].Image = newImageName
	_, err = clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}
