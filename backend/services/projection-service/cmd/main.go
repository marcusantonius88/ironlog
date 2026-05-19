package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcus/ironlog/backend/services/projection-service/internal/application"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	projectionService := application.NewProjectionService()
	_ = projectionService

	// Setup HTTP routes for read models
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/projections/exercise/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "not_implemented"})
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "# Projection metrics\n")
	})

	log.Printf("Projection service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
