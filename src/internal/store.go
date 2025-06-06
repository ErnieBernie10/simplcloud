package internal

import (
	"fmt"
	"io"
	"os"
)

type StoreApp struct {
	Name                  string
	DockerComposeTemplate string
	EnvTemplate           string
	Meta                  AppMeta
}

type IStore interface {
	GetApps() ([]StoreApp, error)
	GetApp(name string) (*StoreApp, error)
}

func (app *StoreApp) Logo(w io.Writer) error {
	storeDir := os.Getenv("STORE_DIR")
	storeFs := os.DirFS(storeDir)
	f, err := storeFs.Open(fmt.Sprintf("%s/%s", app.Name, app.Meta.Logo))
	if err != nil {
		return err
	}
	defer f.Close()
	ff, err := io.ReadAll(f)
	_, err = w.Write(ff)
	if err != nil {
		return err
	}
	return nil
}
