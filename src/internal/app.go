package internal

import "fmt"

type App struct {
	Name          string
	Version       string
	Description   string
	Author        string
	TargetDir     string
	DockerCompose string
	EnvTemplate   string
}

func GetApp(name string) (*App, error) {
	file, err := Opt.ReadFile(fmt.Sprintf("apps/%s/docker-compose.yml", name))
	if err != nil {
		return nil, err
	}
	env, err := Opt.ReadFile(fmt.Sprintf("apps/%s/.env", name))
	if err != nil {
		return nil, err
	}
	return &App{
		Name:          name,
		Version:       "1.0.0",        // todo
		Description:   "A sample app", // todo
		Author:        "John Doe",     // todo
		TargetDir:     fmt.Sprintf("%s/apps/%s", TargetDir, name),
		DockerCompose: string(file),
		EnvTemplate:   string(env),
	}, nil
}
