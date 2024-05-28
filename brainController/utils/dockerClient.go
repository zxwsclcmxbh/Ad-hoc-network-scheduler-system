package utils

import (
	"github.com/docker/docker/client"
)

var DockerClient *client.Client

func DockerClientInit() {

	var err error
	DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err.Error())
		return
	}
}
