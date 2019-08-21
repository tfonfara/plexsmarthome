package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"github.com/tfonfara/plexsmarthome/app"
	"github.com/tfonfara/plexsmarthome/constants"
	"github.com/tfonfara/plexsmarthome/plex"
)

func main() {
	log.Printf("=== Welcome to PlexSmartHome ===\n")
	log.Printf("Container version %s\n", constants.Version)

	eh := app.NewEventHandler()
	wh := plex.NewWebhook(eh.WebhookHandler)

	r := mux.NewRouter()
	r.HandleFunc("/", wh.Handler).Methods(http.MethodPost)
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	port := os.Getenv("PLEX_PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Error starting container: %v\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
