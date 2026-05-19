package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "github.com/marcus/ironlog/backend/services/parser-service/internal/adapters/inbound/http"
	"github.com/marcus/ironlog/backend/services/parser-service/internal/application"
	"github.com/marcus/ironlog/backend/shared/infra"
)

func main() {
	// Get environment configuration
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	if kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize Kafka producer (for event publishing)
	kafkaProducer := infra.NewKafkaProducer(kafkaBrokers, "parsing-events")
	defer kafkaProducer.Close()

	// Initialize application service
	parsingService := application.NewParsingService()

	// Initialize HTTP handler
	handler := httpAdapter.NewParsingHandler(parsingService)

	// Setup HTTP routes
	http.HandleFunc("/parse", handler.ParseDSL)
	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Parser service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// metricsHandler provides Prometheus metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "# HELP parsing_requests_total Total parsing requests\n")
	fmt.Fprint(w, "# TYPE parsing_requests_total counter\n")
	fmt.Fprint(w, "parsing_requests_total 0\n")
}
