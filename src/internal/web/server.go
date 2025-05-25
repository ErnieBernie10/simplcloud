package web

import (
	"log"
	"net/http"

	"github.com/ErnieBernie10/simplecloud/internal/web/controller"
	"github.com/ErnieBernie10/simplecloud/internal/web/core"
)

func Serve() {
	tmplMngr := core.NewTemplateManager("internal/web/templates")

    logger := log.New(log.Writer(), "web: ", log.LstdFlags)

	appContext := &core.AppContext{
		Template: tmplMngr,
		Logger:  logger,
	}

	mux := http.NewServeMux() // Use ExactServeMux to avoid route duplication

	controller.SetupHome(mux, appContext)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
    log.Println("Server running on http://localhost:8080")
	err := server.ListenAndServe() // this blocks and waits for requests
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

