package internal

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func IsMasterRunning() bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}

	projectName := filepath.Base(TargetDir)

	args := filters.NewArgs()
	args.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	args.Add("status", "running")

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		return false
	}

	return len(containers) == 2 // traefik and whoami

}

func RunMaster() error {
	cmd := exec.Command("docker", "compose", "-f", fmt.Sprintf("%s/docker-compose.yml", TargetDir), "up", "-d")
	return cmd.Run()
}

func StopMaster() error {
	cmd := exec.Command("docker", "compose", "-f", fmt.Sprintf("%s/docker-compose.yml", TargetDir), "down")
	return cmd.Run()
}

func RestartMaster() error {
	cmd := exec.Command("docker", "compose", "-f", fmt.Sprintf("%s/docker-compose.yml", TargetDir), "restart")
	return cmd.Run()
}
