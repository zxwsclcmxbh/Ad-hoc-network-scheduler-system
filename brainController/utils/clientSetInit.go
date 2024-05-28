package utils

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var ClientSet *kubernetes.Clientset
var MetricsClientSet *metricsv.Clientset

func ClientSetInit() {

	var kubeConfig *string

	// home是家目录，如果能取得家目录的值，就可以用来做默认值
	if home := homedir.HomeDir(); home != "" {
		// 如果输入了kubeConfig参数，该参数的值就是kubeConfig文件的绝对路径，
		// 如果没有输入kubeConfig参数，就用默认路径~/.kube/config
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		// 如果取不到当前用户的家目录，就没办法设置kubeConfig的默认目录了，只能从入参中取
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()
	var err error
	// 从本机加载kubeConfig配置文件，因此第一个参数为空字符串
	ConfigString, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)

	// kubeConfig加载失败就直接退出了
	if err != nil {
		panic(err.Error())
	}

	// 实例化clientSet对象
	ClientSet, err = kubernetes.NewForConfig(ConfigString)
	MetricsClientSet, err = metricsv.NewForConfig(ConfigString)
	if err != nil {
		panic(err.Error())
		return
	}
}
