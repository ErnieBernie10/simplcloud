package service

import (
	"encoding/json"
	"errors"
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"io"
	"net/http"
)

type Store struct {
	BaseURL string
}

func NewStoreService(baseURL string) internal.IStore {
	return &Store{
		BaseURL: baseURL,
	}
}

func (s *Store) GetApps() ([]internal.StoreApp, error) {
	resp, err := http.Get(s.BaseURL + "/store")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(bodyBytes))
	}

	var apps []internal.StoreApp
	err = json.NewDecoder(resp.Body).Decode(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Store) GetApp(name string) (*internal.StoreApp, error) {
	resp, err := http.Get(s.BaseURL + "/app/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(bodyBytes))
	}

	var app internal.StoreApp
	err = json.NewDecoder(resp.Body).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
