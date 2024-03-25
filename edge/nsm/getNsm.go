package nsm

import (
	"bufio"
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	cli "github.com/docker/docker/client"
)

func GetNsm(containerID string) ([]string, error) {
	var res []string
	// 创建 Docker 客户端
	cli, err := cli.NewClientWithOpts(cli.WithAPIVersionNegotiation())
	if err != nil {
		return res, err
	}
	// 构建执行命令所需的参数
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"vtysh", "-c", "show ip ospf neighbor"},
	}
	// 在容器中执行命令
	resp, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return res, err
	}
	// 获取命令执行的输出
	respHijacked, err := cli.ContainerExecAttach(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return res, err
	}
	defer respHijacked.Close()
	// 从输出中读取并处理结果

	scanner := bufio.NewScanner(respHijacked.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过标题行
		if strings.Contains(line, "Neighbor ID") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 6 {
			res = append(res, fields[5])
		}
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}
