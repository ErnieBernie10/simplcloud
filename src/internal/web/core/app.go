package core

import (
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"log"
)

const (
	BaseLayout = "layout"
)

type AppContext struct {
	Template     *TemplateManager
	Logger       *log.Logger
	StoreService internal.IStore
}
