package internal

import (
	"embed"
	"os"
)

var Opt embed.FS

var TargetDir string

func Load() {
	TargetDir = os.Getenv("TARGET_DIR")
}
