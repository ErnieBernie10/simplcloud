package core

import (
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"log"
)

type AppContext struct {
	Logger       *log.Logger
	StoreDir     string
	StoreService internal.IStore
}

func NewAppContext(logger *log.Logger, storeDir string, storeService internal.IStore) *AppContext {

	appContext := &AppContext{
		Logger:       logger,
		StoreDir:     storeDir,
		StoreService: storeService,
	}

	return appContext
}
