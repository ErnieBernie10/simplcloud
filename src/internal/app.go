package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"html/template"
	"os"
	"os/exec"
	"path"
)

type App struct {
	Context               *RunContext
	Name                  string
	Version               string
	Description           string
	Author                string
	TargetDir             string
	DockerComposeTemplate string
	EnvTemplate           string
	DockerComposeFile     string
	EnvFile               string
	Meta                  AppMeta
}

type AppMetaFile struct {
	App AppMeta `yaml:"app"`
}

type AppMeta struct {
	DisplayName string `yaml:"display_name"`
	Logo        string `yaml:"logo"`
	Description string `yaml:"description"`
}

func (r *RunContext) GetApp(name string) (*App, error) {
	file, err := Opt.ReadFile(fmt.Sprintf("apps/%s/docker-compose.yml", name))
	if err != nil {
		return nil, err
	}
	metaFile, err := Opt.ReadFile(fmt.Sprintf("apps/%s/meta.yml", name))
	if err != nil {
		return nil, err
	}
	env, err := Opt.ReadFile(fmt.Sprintf("apps/%s/.env", name))
	if err != nil {
		return nil, err
	}
	appTargetDir := fmt.Sprintf("%s/apps/%s", r.TargetDir, name)

	var meta AppMetaFile
	err = yaml.Unmarshal(metaFile, &meta)
	if err != nil {
		return nil, err
	}
	return &App{
		Context:               r,
		Name:                  name,
		Version:               "1.0.0",        // todo
		Description:           "A sample app", // todo
		Author:                "John Doe",     // todo
		TargetDir:             appTargetDir,
		DockerComposeTemplate: string(file),
		EnvTemplate:           string(env),
		DockerComposeFile:     fmt.Sprintf("%s/docker-compose.yml", appTargetDir),
		EnvFile:               fmt.Sprintf("%s/.env", appTargetDir),
		Meta:                  meta.App,
	}, nil
}

func (app *App) Deploy() error {
	if err := os.MkdirAll(path.Dir(app.DockerComposeFile), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(app.DockerComposeFile, []byte(app.DockerComposeTemplate), 0644); err != nil {
		return err
	}
	tmpl, err := template.New("env").Parse(app.EnvTemplate)
	if err != nil {
		return err
	}
	config, err := app.Context.config()
	if err != nil {
		return err
	}

	envFile, err := os.Create(app.EnvFile)
	if err != nil {
		return err
	}
	defer envFile.Close()
	// TODO: Accept subdomain as parameter
	// TODO: Accept port as parameter
	err = tmpl.Execute(envFile, map[string]string{
		"Domain": config.Config.Domain,
		"Email":  config.Config.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	cmd := exec.Command("docker", "compose", "-f", app.DockerComposeFile, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
