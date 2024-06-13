package utils

import (
	_ "flag"
	"log"
	_ "path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	_ "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	_ "k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/util/homedir"
)

var ConfigString *rest.Config
var OssClient *minio.Client

func OSSInit() {
	var err error
	OssClient, err = minio.New(Config.Minio.URL, &minio.Options{Creds: credentials.NewStaticV4(Config.Minio.Username, Config.Minio.Password, "")})
	if err != nil {
		log.Fatalln(err)
	}
}
