package dockerService

import "cloud/brainController/service/shellService"

func DeleteImageInRegistry(image string) error {
	cmd := "docker exec registry rm -rf /var/lib/registry/docker/registry/v2/repositories/" + image
	err, _ := shellService.ExecShell(cmd)
	if err != nil {
		return err
	}
	return nil
}
