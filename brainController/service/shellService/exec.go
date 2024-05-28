package shellService

import (
	"bytes"
	"log"
	"os/exec"
)

func ExecShell(cmdString string) (error, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Printf("cmd.Run() failed with %s\n\n", err)
		return err, stdout.String()
	}
	return nil, stdout.String()
}
