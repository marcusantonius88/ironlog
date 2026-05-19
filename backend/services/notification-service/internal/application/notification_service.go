package application

import (
	"log"

	"github.com/marcus/ironlog/backend/services/notification-service/internal/domain"
)

// NotificationService handles notification dispatch
type NotificationService struct{}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendPersonalRecordNotification sends a PR notification
func (ns *NotificationService) SendPersonalRecordNotification(userID, exerciseName string, value float64, recordType string) error {
	message := exerciseName + " " + recordType + ": " + string(rune(int(value)))
	log.Printf("[NOTIFICATION] PR for user %s: %s", userID, message)
	return nil
}

// SendProgressionNotification sends a progression notification
func (ns *NotificationService) SendProgressionNotification(userID, exerciseName string) error {
	message := "Progress detected in " + exerciseName
	log.Printf("[NOTIFICATION] Progression for user %s: %s", userID, message)
	return nil
}

// SendRegressionNotification sends a regression notification
func (ns *NotificationService) SendRegressionNotification(userID, exerciseName string) error {
	message := "Regression detected in " + exerciseName + ". Consider deloading."
	log.Printf("[NOTIFICATION] Regression for user %s: %s", userID, message)
	return nil
}

// SendLoadIncreaseRecommendation sends a load increase suggestion
func (ns *NotificationService) SendLoadIncreaseRecommendation(userID, exerciseName string, suggestedWeight float64) error {
	message := "Try increasing " + exerciseName + " weight"
	log.Printf("[NOTIFICATION] Recommendation for user %s: %s (suggested: %f)", userID, message, suggestedWeight)
	return nil
}

// RouteNotification routes a notification to the appropriate channel
func (ns *NotificationService) RouteNotification(notification *domain.Notification) error {
	switch notification.Channel {
	case "LOG":
		log.Printf("[%s] %s: %s", notification.Type, notification.UserID, notification.Message)
	case "EMAIL":
		// Email sending would be implemented here
		log.Printf("[EMAIL] To: %s, Subject: %s", notification.UserID, notification.Title)
	case "TELEGRAM":
		// Telegram sending would be implemented here
		log.Printf("[TELEGRAM] To: %s, Message: %s", notification.UserID, notification.Message)
	}
	return nil
}
