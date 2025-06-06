package web

import (
	"errors"
	"fmt"
	"github.com/ErnieBernie10/simplecloud/src/internal/web/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ErnieBernie10/simplecloud/src/internal/web/controller"
	"github.com/ErnieBernie10/simplecloud/src/internal/web/core"
)

func Serve() {
	root, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	root = filepath.Dir(root)
	fmt.Println("root:", root)
	tmplMngr := core.NewTemplateManager(filepath.Join(root, "templates"))

	logger := log.New(log.Writer(), "web: ", log.LstdFlags)

	mux := http.NewServeMux() // Use ExactServeMux to avoid route duplication

	appContext := &core.AppContext{
		Template:     tmplMngr,
		Logger:       logger,
		StoreService: service.NewStoreService(os.Getenv("BASE_URL")),
	}

	staticDir := filepath.Join(root, "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	controller.SetupHome(mux, appContext)
	controller.SetupStore(mux, appContext)

	server := &http.Server{
		Addr:    ":8181",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop // Wait for a termination signal
		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error shutting down server: %s", err)
		}
	}()

	log.Println("Starting server on :8181")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %s", err)
	}
}
