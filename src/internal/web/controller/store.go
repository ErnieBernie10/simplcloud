package controller

import (
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"github.com/ErnieBernie10/simplecloud/src/internal/web/core"
	"net/http"
)

type StoreController struct {
	AppContext *core.AppContext
}

func NewStoreController(appContext *core.AppContext) *StoreController {
	return &StoreController{
		AppContext: appContext,
	}
}

type StoreVM struct {
	Apps []internal.StoreApp
	Cdn  string
}

func (c *StoreController) Store(w http.ResponseWriter, r *http.Request) {
	apps, err := c.AppContext.StoreService.GetApps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vm := StoreVM{
		Apps: apps,
		Cdn:  c.AppContext.Cdn,
	}
	c.AppContext.Template.Render(w, "store.html", vm, core.BaseLayout)
}

func SetupStore(mux *http.ServeMux, context *core.AppContext) {
	storeController := NewStoreController(context)

	// Register routes
	mux.HandleFunc("/store", storeController.Store)
}
