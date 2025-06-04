package core

import "github.com/ErnieBernie10/simplecloud/src/internal"

type IStore interface {
	GetApp(name string) (*internal.StoreApp, error)
}
