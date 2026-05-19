package domain

// Notification represents a notification to be sent
type Notification struct {
	ID      string
	UserID  string
	Type    string // PR, PROGRESSION, REGRESSION, SUGGESTION
	Title   string
	Message string
	Channel string // LOG, EMAIL, TELEGRAM
	SentAt  string
	ReadAt  string
}

// NotificationEvent represents notification system event
type NotificationEvent struct {
	EventID string
	Type    string
	UserID  string
	Payload interface{}
}
