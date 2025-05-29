package core

import (
	"log"
	"net/http"

	"github.com/ErnieBernie10/simplecloud/src/internal"
)

const (
	BaseLayout = "layout"
)

type AppContext struct {
	Template   *TemplateManager
	Logger     *log.Logger
	RunContext *internal.RunContext
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
