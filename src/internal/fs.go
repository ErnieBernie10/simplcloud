package internal

import (
	"embed"
	"os"
)

//go:embed apps/**/docker-compose.yml apps/**/.env traefik/.env traefik/docker-compose.yml
var Opt embed.FS

var TargetDir string


func Load() {
	TargetDir = os.Getenv("TARGET_DIR")
}
