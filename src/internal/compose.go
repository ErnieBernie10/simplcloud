package internal

import (
	"gopkg.in/yaml.v3"
)

type ComposeFile struct {
	Services map[string]struct {
		Image string `yaml:"image"`
	} `yaml:"services"`
}

func ParseCompose(compose string) (ComposeFile, error) {
	var file ComposeFile
	if err := yaml.Unmarshal([]byte(compose), &file); err != nil {
		return file, err
	}

	return file, nil
}
