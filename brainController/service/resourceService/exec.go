package resourceService

import (
	"bytes"
	"cloud/brainController/utils"
	"context"
	"errors"
	"io"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecCommandAndGetResult(clientSet *kubernetes.Clientset, namespace string, podName string, container string, command []string) (string, error) {
	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   command,
			// Command: []string{"sleep", "15"},
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			TTY:    false,
		}, scheme.ParameterCodec)
	exec, _ := remotecommand.NewSPDYExecutor(utils.ConfigString, "POST", req.URL())
	resultChan := make(chan string)
	errChan := make(chan error)
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	go func() {
		buffer := new(bytes.Buffer)
		if err := exec.Stream(remotecommand.StreamOptions{Stdout: buffer, Stdin: os.Stdin, Stderr: os.Stderr, Tty: false}); err != nil {
			errChan <- err
		}
		resultChan <- buffer.String()
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", errors.New("deadline exceeded")
	}
}

func ExecCommandAndCopy(clientSet *kubernetes.Clientset, namespace string, podName string, container string, fileName string, writer *io.PipeWriter) {
	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   []string{"tar", "cf", "-", "/home/jovyan/work/" + fileName},
			Container: container,
			// Command: []string{"sleep", "15"},
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			TTY:    false,
		}, scheme.ParameterCodec)
	exec, _ := remotecommand.NewSPDYExecutor(utils.ConfigString, "POST", req.URL())
	go func() {
		defer writer.Close()
		err := exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Stdout: writer,
			Tty:    false,
		})
		if err != nil {
			writer.CloseWithError(err)
			return
		}
	}()
	return
}
