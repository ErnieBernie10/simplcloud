package controller

import (
	"net/http"

	"github.com/ErnieBernie10/simplecloud/src/internal/web/core"
)

type HomeController struct {
	AppContext *core.AppContext
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	c.AppContext.Template.Render(w, "index.html", nil, core.BaseLayout)
}

func NewHomeController(appContext *core.AppContext) *HomeController {
	return &HomeController{
		AppContext: appContext,
	}
}

func SetupHome(mux *http.ServeMux, context *core.AppContext) {
	homeController := NewHomeController(context)

	// Register routes
	mux.HandleFunc("/", core.ExactRoute("/", homeController.Index))
}
