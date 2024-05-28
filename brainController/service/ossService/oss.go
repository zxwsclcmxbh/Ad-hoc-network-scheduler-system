package ossService

import (
	"cloud/brainController/utils"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func Upload(filename string, file *io.PipeReader) (string, error) {
	client := utils.OssClient
	_, err := client.PutObject(context.TODO(), "phm-data", filename, file, -1, minio.PutObjectOptions{})
	// _, err := client.Object.Put(context.Background(), filename, file, nil)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("http://%s/%s/%s", utils.Config.Minio.URL, "phm-data", filename)
	return result, nil
	// result := client.Object.GetObjectURL(filename)
	// return result.String(), nil
}
