package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ErnieBernie10/simplecloud/src/internal/web/controller"
	"github.com/ErnieBernie10/simplecloud/src/internal/web/core"
)

func Serve() {
	root, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	root = filepath.Dir(root)
	tmplMngr := core.NewTemplateManager(filepath.Join(root, "templates"))

	logger := log.New(log.Writer(), "web: ", log.LstdFlags)

	appContext := &core.AppContext{
		Template: tmplMngr,
		Logger:   logger,
	}

	mux := http.NewServeMux() // Use ExactServeMux to avoid route duplication

	controller.SetupHome(mux, appContext)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	staticDir := filepath.Join(root, "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("static/", http.StripPrefix("static/", fs))

	log.Println("Server running on http://localhost:8080")
	err = server.ListenAndServe() // this blocks and waits for requests
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
