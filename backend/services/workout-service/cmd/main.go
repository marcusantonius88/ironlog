package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "github.com/marcus/ironlog/backend/services/workout-service/internal/adapters/inbound/http"
	"github.com/marcus/ironlog/backend/services/workout-service/internal/adapters/outbound/postgres"
	"github.com/marcus/ironlog/backend/services/workout-service/internal/application"
	"github.com/marcus/ironlog/backend/shared/infra"
)

func main() {
	// Get configuration from environment
	dbConfig := infra.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	if dbConfig.Host == "" {
		dbConfig.Host = "localhost"
	}
	if dbConfig.User == "" {
		dbConfig.User = "postgres"
	}
	if dbConfig.Password == "" {
		dbConfig.Password = "password"
	}
	if dbConfig.Database == "" {
		dbConfig.Database = "ironlog"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	db, err := infra.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize event store
	eventStore := postgres.NewPostgresEventStore(db)

	// Initialize application service
	// Note: In a real implementation, we'd also need a workout repository
	workoutService := application.NewWorkoutService(nil, eventStore)

	// Initialize HTTP handler
	handler := httpAdapter.NewWorkoutHandler(workoutService)

	// Setup HTTP routes
	http.HandleFunc("/workouts/start", handler.StartWorkout)
	http.HandleFunc("/workouts/sets/perform", handler.PerformSet)
	http.HandleFunc("/workouts/finish", handler.FinishWorkout)
	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Workout service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// metricsHandler provides Prometheus metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "# HELP events_processed_total Total events processed\n")
	fmt.Fprint(w, "# TYPE events_processed_total counter\n")
	fmt.Fprint(w, "events_processed_total 0\n")
}
