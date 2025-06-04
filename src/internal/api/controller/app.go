package controller

import (
	"errors"
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"github.com/ErnieBernie10/simplecloud/src/internal/api/core"
	"net/http"
	"strings"
)

type AppController struct {
	AppContext *core.AppContext
}

func NewAppController(appContext *core.AppContext) *AppController {
	return &AppController{
		AppContext: appContext,
	}
}

func SetupApp(mux *http.ServeMux, context *core.AppContext) {
	appController := NewAppController(context)
	// Register routes
	mux.HandleFunc("GET /app/{app}/logo", appController.Logo)
	mux.HandleFunc("GET /app/{app}", appController.GetApp)
}

func (c *AppController) GetApp(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("app")

	app, err := c.AppContext.StoreService.GetApp(appName)
	if err != nil {
		if errors.Is(err, internal.ErrAppNotFound) {
			http.Error(w, "App not found", http.StatusNotFound)
		}
	}

	err = core.WriteJSON(w, app, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AppController) Logo(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("app")
	if appName == "" {
		http.Error(w, "Invalid app name", http.StatusBadRequest)
		return
	}

	app, err := c.AppContext.StoreService.GetApp(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if strings.HasSuffix(app.Meta.Logo, ".svg") {
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	if strings.HasSuffix(app.Meta.Logo, ".png") {
		w.Header().Set("Content-Type", "image/png")
	}
	err = app.Logo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
