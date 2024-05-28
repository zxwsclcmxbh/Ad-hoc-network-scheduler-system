package dockerService

import (
	"cloud/brainController/utils"
	"context"
	"github.com/docker/docker/api/types"
)

func ImageList() ([]types.ImageSummary, error) {
	dockerClient := utils.DockerClient
	return dockerClient.ImageList(context.Background(), types.ImageListOptions{})
}
