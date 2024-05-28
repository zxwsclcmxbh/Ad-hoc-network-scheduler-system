package dockerService

import (
	"cloud/brainController/utils"
	"context"
	"encoding/base64"
	"github.com/docker/docker/api/types"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
)

func PushImage(imageName string) (string, error) {
	dockerClient := utils.DockerClient
	user := ""
	password := ""
	authConfig := types.AuthConfig{Username: user, Password: password}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, ImagePushErr := dockerClient.ImagePush(context.TODO(), imageName, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	if ImagePushErr != nil {
		return "", ImagePushErr
	}

	body, ReadAllErr := io.ReadAll(out)
	if ReadAllErr != nil {
		return "", ReadAllErr
	}

	return string(body), nil
}
