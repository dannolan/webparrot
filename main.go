package main

import (
	"log"
	"net/http"
	"os"

	"github.com/webparrot/handlers"
)

func main() {
	logger := log.New(os.Stdout, "webparrot:", log.LstdFlags)
	logger.Println("Loading...")
	listenAddr := ":5000"
	router := http.NewServeMux()
	router.Handle("/healthz", handlers.HealthHandler{
		Logger: logger,
	})
	router.Handle("/api/v1/parrot", handlers.ParrotHandler{
		Logger: logger,
	})

	log.Fatal(http.ListenAndServe(listenAddr, router))
}
