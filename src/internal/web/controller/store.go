package controller

import (
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

func (c *StoreController) Store(w http.ResponseWriter, r *http.Request) {
	apps, err := c.AppContext.RunContext.GetApps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.AppContext.Template.Render(w, "store.html", apps, core.BaseLayout)
}

func SetupStore(mux *http.ServeMux, context *core.AppContext) {
	storeController := NewStoreController(context)

	// Register routes
	mux.HandleFunc("/store", core.ExactRoute("/store", storeController.Store))
}
