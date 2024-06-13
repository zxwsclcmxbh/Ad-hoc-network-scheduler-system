package utils

// docker buildx build http://192.168.1.3:9000/phm-model-resource/model-635-17.tar -t 127.0.0.1:5000/models-635-17:v0  --push --platform="linux/amd64,linux/arm64,linux/arm/v7"
import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

func RunDockerBuildX(url string, image string) error {
	cmd := exec.Command("docker", "buildx", "build", url, "-t", image, "--output=type=image,push=true,registry.insecure=true", "--platform=linux/amd64,linux/arm64,linux/arm/v7")
	stderr, _ := cmd.StderrPipe()
	log.Println(cmd.Start())
	io.Copy(os.Stdout, stderr)
	err := cmd.Wait()
	if Err, ok := err.(*exec.ExitError); ok {
		return errors.New(Err.Error())
	}
	return err
}
