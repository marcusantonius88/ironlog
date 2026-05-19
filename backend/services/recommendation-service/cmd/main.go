package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcus/ironlog/backend/services/recommendation-service/internal/application"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	recommendationService := application.NewRecommendationService()
	_ = recommendationService

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "# HELP load_increase_suggestions_total Total load increase suggestions\n")
		fmt.Fprint(w, "# TYPE load_increase_suggestions_total counter\n")
		fmt.Fprint(w, "load_increase_suggestions_total 0\n")
	})

	log.Printf("Recommendation service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
