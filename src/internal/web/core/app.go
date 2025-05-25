package core

import (
	"log"
	"net/http"
)

const (
	BaseLayout = "layout"
)

type AppContext struct {
	Template *TemplateManager
	Logger  *log.Logger
}

func ExactRoute(route string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == route {
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
