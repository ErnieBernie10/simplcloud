package controller

import (
	"github.com/ErnieBernie10/simplecloud/src/internal/api/core"
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

}

func SetupStore(mux *http.ServeMux, context *core.AppContext) {
	storeController := NewStoreController(context)

	// Register routes
	mux.HandleFunc("GET /store", storeController.Store)
}
