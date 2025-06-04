package api

import (
	"errors"
	"fmt"
	"github.com/ErnieBernie10/simplecloud/src/internal/api/controller"
	"github.com/ErnieBernie10/simplecloud/src/internal/api/core"
	"github.com/ErnieBernie10/simplecloud/src/internal/api/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func Serve() {
	root, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return
	}
	root = filepath.Dir(root)
	fmt.Println("root:", root)

	logger := log.New(log.Writer(), "web: ", log.LstdFlags)

	mux := http.NewServeMux() // Use ExactServeMux to avoid route duplication

	storeDir := os.Getenv("STORE_DIR")
	storeService, err := services.NewStoreService(storeDir + "/store.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	appContext := core.NewAppContext(logger, storeDir+"/store.json", storeService)

	controller.SetupStore(mux, appContext)
	controller.SetupApp(mux, appContext)

	server := &http.Server{
		Addr:    ":8080",
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

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %s", err)
	}
}
