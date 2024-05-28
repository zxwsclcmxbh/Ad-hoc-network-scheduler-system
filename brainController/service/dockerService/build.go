package dockerService

import (
	"cloud/brainController/utils"
	"context"
	"github.com/docker/docker/api/types"
	"io"
)

// BuildImage 构建镜像
func BuildImage(tarUrl string, tags string) (string, error) {
	dockerClient := utils.DockerClient
	options := types.ImageBuildOptions{
		Tags:          []string{tags},
		RemoteContext: tarUrl,
	}
	rsp, imageBuildErr := dockerClient.ImageBuild(context.Background(), nil, options)
	if imageBuildErr != nil {
		return "", imageBuildErr
	}
	body, ReadAllErr := io.ReadAll(rsp.Body)
	if ReadAllErr != nil {
		return "", ReadAllErr
	}
	return string(body), nil
}
