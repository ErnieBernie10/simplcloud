package internal

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type RunContext struct {
	TargetDir         string
	DockerComposeFile string
	Context           context.Context
}

type Config struct {
	Domain string `toml:"domain"`
	Email  string `toml:"email"`
}
type TomlConfig struct {
	Config Config `toml:"Config"`
}

func NewRunContext(targetDir string, c context.Context) *RunContext {
	return &RunContext{
		TargetDir:         targetDir,
		DockerComposeFile: fmt.Sprintf("%s/docker-compose.yml", targetDir),
		Context:           c,
	}
}

func (c *RunContext) DockerProjectName() string {
	return strings.ToLower(filepath.Base(c.TargetDir))
}

func (c *RunContext) IsMasterRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}

	projectName := c.DockerProjectName()

	args := filters.NewArgs()
	args.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	args.Add("status", "running")

	containers, err := cli.ContainerList(c.Context, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		fmt.Println("Error listing containers:", err)
		return false
	}

	return len(containers) == 2 // traefik and whoami

}

func (c *RunContext) Bootstrap(config TomlConfig) error {
	content, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(c.TargetDir, 0755); err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/config.toml", c.TargetDir), content, 0644)
	if err != nil {
		return err
	}

	if err = c.createInitDockerCompose(); err != nil {
		return err
	}
	if err = c.createInitEnv(config); err != nil {
		return err
	}

	return nil
}

func (c *RunContext) createInitEnv(config TomlConfig) error {
	env, err := Opt.ReadFile("traefik/.env")
	if err != nil {
		return err
	}
	tmpl, err := template.New("env").Parse(string(env))
	if err != nil {
		return err
	}
	envFile, err := os.Create(fmt.Sprintf("%s/.env", c.TargetDir))
	if err != nil {
		return err
	}
	defer envFile.Close()
	err = tmpl.Execute(envFile, map[string]string{
		"Domain": config.Config.Domain,
		"Email":  config.Config.Email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RunContext) createInitDockerCompose() error {
	traefik, err := Opt.ReadFile("traefik/docker-compose.yml")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/docker-compose.yml", c.TargetDir), traefik, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *RunContext) RunMaster() error {
	cmd := exec.Command("docker", "compose", "-f", c.DockerComposeFile, "up", "-d")
	return cmd.Run()
}

func (c *RunContext) StopMaster() error {
	cmd := exec.Command("docker", "compose", "-f", c.DockerComposeFile, "down")
	return cmd.Run()
}

func (c *RunContext) RestartMaster() error {
	cmd := exec.Command("docker", "compose", "-f", c.DockerComposeFile, "restart")
	return cmd.Run()
}

func (c *RunContext) config() (TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", c.TargetDir), &config); err != nil {
		return config, err
	}
	return config, nil
}
