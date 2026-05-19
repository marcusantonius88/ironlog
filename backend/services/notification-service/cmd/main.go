package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marcus/ironlog/backend/services/notification-service/internal/application"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	notificationService := application.NewNotificationService()
	_ = notificationService

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "# HELP notifications_sent_total Total notifications sent\n")
		fmt.Fprint(w, "# TYPE notifications_sent_total counter\n")
		fmt.Fprint(w, "notifications_sent_total 0\n")
	})

	log.Printf("Notification service starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
