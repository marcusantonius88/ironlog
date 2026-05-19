package events

import (
	"time"

	"github.com/google/uuid"
)

// EventEnvelope defines the contract for all domain events
type EventEnvelope struct {
	EventID       string            `json:"event_id"`
	EventType     string            `json:"event_type"`
	AggregateID   string            `json:"aggregate_id"`
	CorrelationID string            `json:"correlation_id"`
	Payload       interface{}       `json:"payload"`
	CreatedAt     time.Time         `json:"created_at"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// NewEventEnvelope creates a new event envelope
func NewEventEnvelope(eventType string, aggregateID string, correlationID string, payload interface{}) *EventEnvelope {
	return &EventEnvelope{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		AggregateID:   aggregateID,
		CorrelationID: correlationID,
		Payload:       payload,
		CreatedAt:     time.Now().UTC(),
		Metadata:      make(map[string]string),
	}
}

// Event types constants
const (
	// Workout events
	WorkoutStartedEvent   = "WorkoutStarted"
	WorkoutFinishedEvent  = "WorkoutFinished"
	ExerciseStartedEvent  = "ExerciseStarted"
	ExerciseFinishedEvent = "ExerciseFinished"
	SetPerformedEvent     = "SetPerformed"
	WorkoutNoteAddedEvent = "WorkoutNoteAdded"

	// Parser events
	ParsingRequestedEvent = "ParsingRequested"

	// Analytics events
	PerformanceImprovedEvent    = "PerformanceImproved"
	PerformanceRegressedEvent   = "PerformanceRegressed"
	PersonalRecordAchievedEvent = "PersonalRecordAchieved"
	VolumeIncreasedEvent        = "VolumeIncreased"
	LoadIncreaseSuggestedEvent  = "LoadIncreaseSuggested"

	// Projection events
	ProjectionUpdatedEvent = "ProjectionUpdated"

	// Recommendation events
	LoadIncreaseRecommendedEvent = "LoadIncreaseRecommended"
	DeloadRecommendedEvent       = "DeloadRecommended"

	// Notification events
	NotificationSentEvent = "NotificationSent"
)
