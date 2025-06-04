package core

import (
	"log"
)

type AppContext struct {
	Logger       *log.Logger
	StoreDir     string
	StoreService IStore
}

func NewAppContext(logger *log.Logger, storeDir string, storeService IStore) *AppContext {

	appContext := &AppContext{
		Logger:       logger,
		StoreDir:     storeDir,
		StoreService: storeService,
	}

	return appContext
}
