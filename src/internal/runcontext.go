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

func (r *RunContext) DockerProjectName() string {
	return strings.ToLower(filepath.Base(r.TargetDir))
}

func (r *RunContext) IsMasterRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}

	projectName := r.DockerProjectName()

	args := filters.NewArgs()
	args.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	args.Add("status", "running")

	containers, err := cli.ContainerList(r.Context, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		fmt.Println("Error listing containers:", err)
		return false
	}

	return len(containers) == 2 // traefik and whoami

}

func (r *RunContext) Bootstrap(config TomlConfig) error {
	content, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(r.TargetDir, 0755); err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/config.toml", r.TargetDir), content, 0644)
	if err != nil {
		return err
	}

	if err = r.createInitDockerCompose(); err != nil {
		return err
	}
	if err = r.createInitEnv(config); err != nil {
		return err
	}

	return nil
}

func (r *RunContext) createInitEnv(config TomlConfig) error {
	env, err := Opt.ReadFile("traefik/.env")
	if err != nil {
		return err
	}
	tmpl, err := template.New("env").Parse(string(env))
	if err != nil {
		return err
	}
	envFile, err := os.Create(fmt.Sprintf("%s/.env", r.TargetDir))
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

func (r *RunContext) createInitDockerCompose() error {
	traefik, err := Opt.ReadFile("traefik/docker-compose.yml")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/docker-compose.yml", r.TargetDir), traefik, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *RunContext) RunMaster() error {
	cmd := exec.Command("docker", "compose", "-f", r.DockerComposeFile, "up", "-d")
	return cmd.Run()
}

func (r *RunContext) StopMaster() error {
	cmd := exec.Command("docker", "compose", "-f", r.DockerComposeFile, "down")
	return cmd.Run()
}

func (r *RunContext) RestartMaster() error {
	cmd := exec.Command("docker", "compose", "-f", r.DockerComposeFile, "restart")
	return cmd.Run()
}

func (r *RunContext) config() (TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", r.TargetDir), &config); err != nil {
		return config, err
	}
	return config, nil
}

func (r *RunContext) GetApps() ([]App, error) {

	var apps []App
	files, err := Opt.ReadDir("apps")
	if err != nil {
		return apps, err
	}
	for _, file := range files {
		if file.IsDir() {
			app, err := r.GetApp(file.Name())
			if err != nil {
				return apps, err
			}
			apps = append(apps, *app)
		}
	}

	return apps, nil
}
