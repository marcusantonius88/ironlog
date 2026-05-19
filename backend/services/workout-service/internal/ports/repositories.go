package ports

import (
	"github.com/marcus/ironlog/backend/services/workout-service/internal/domain"
)

// WorkoutRepository defines storage interface for workouts
type WorkoutRepository interface {
	GetByID(id string) (*domain.Workout, error)
	Save(workout *domain.Workout) error
}

// EventStore defines the interface for event persistence (Event Sourcing)
type EventStore interface {
	SaveEvents(aggregateID string, events []interface{}) error
	LoadEvents(aggregateID string) ([]interface{}, error)
	GetEventStream(aggregateID string) ([]interface{}, error)
}

// OutboxRepository defines the interface for outbox pattern
type OutboxRepository interface {
	SaveOutboxEvent(aggregateID string, event interface{}) error
	GetUnpublishedEvents() ([]interface{}, error)
	MarkEventAsPublished(eventID string) error
}
