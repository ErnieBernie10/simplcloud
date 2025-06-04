package main

import (
	"encoding/json"
	"fmt"
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"os"
)

func GetApps(appFs fs.FS, dir string) ([]internal.StoreApp, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var apps []internal.StoreApp
	for _, file := range files {
		if file.IsDir() {
			app, err := GetApp(file.Name(), appFs)
			if err != nil {
				return apps, err
			}
			apps = append(apps, *app)
		}
	}

	return apps, nil

}

func ReadAppFile(appFs fs.FS, name string, filename string) ([]byte, error) {
	file, err := appFs.Open(fmt.Sprintf("%s/%s", name, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetApp(name string, appFs fs.FS) (*internal.StoreApp, error) {
	file, err := ReadAppFile(appFs, name, "docker-compose.yml")
	if err != nil {
		return nil, err
	}
	metaFile, err := ReadAppFile(appFs, name, "meta.yml")
	if err != nil {
		return nil, err
	}
	envFile, err := ReadAppFile(appFs, name, ".env")
	if err != nil {
		return nil, err
	}

	var meta internal.AppMetaFile
	err = yaml.Unmarshal(metaFile, &meta)
	if err != nil {
		return nil, err
	}
	return &internal.StoreApp{
		Name:                  name,
		DockerComposeTemplate: string(file),
		EnvTemplate:           string(envFile),
		Meta:                  meta.App,
	}, nil
}

func buildStore(appFs fs.FS, dir string) {
	log.Println("Building store")
	apps, err := GetApps(appFs, dir)
	if err != nil {
		log.Fatal(err)
	}
	jsonApps, err := json.Marshal(apps)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s/store.json", dir), jsonApps, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Store built")
}
