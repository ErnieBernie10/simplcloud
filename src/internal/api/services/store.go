package services

import (
	"encoding/json"
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"os"
)

type Store struct {
	Apps []internal.StoreApp
}

func NewStoreService(storeDir string) (*Store, error) {
	file, err := os.ReadFile(storeDir)
	if err != nil {
		return nil, err
	}
	var Apps []internal.StoreApp
	err = json.Unmarshal(file, &Apps)
	if err != nil {
		return nil, err
	}
	return &Store{
		Apps: Apps,
	}, nil
}

func (s *Store) GetApp(name string) (*internal.StoreApp, error) {
	for _, app := range s.Apps {
		if app.Name == name {
			return &app, nil
		}
	}
	return nil, internal.ErrAppNotFound
}

func (s *Store) GetApps() ([]internal.StoreApp, error) {
	return s.Apps, nil
}
