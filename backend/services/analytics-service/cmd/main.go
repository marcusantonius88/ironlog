package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcus/ironlog/backend/services/analytics-service/internal/application"
)

func main() {
	// Get configuration
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	if kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Initialize services
	analyticsService := application.NewAnalyticsService()

	// Setup HTTP routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "# HELP personal_records_total Total personal records\n")
		fmt.Fprint(w, "# TYPE personal_records_total counter\n")
		fmt.Fprint(w, "personal_records_total 0\n")
	})

	// Initialize Kafka consumer for events
	_ = kafkaBrokers
	_ = analyticsService

	log.Printf("Analytics service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
