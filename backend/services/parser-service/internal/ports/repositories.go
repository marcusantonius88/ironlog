package ports

import (
	"github.com/marcus/ironlog/backend/services/parser-service/internal/domain"
)

// ParsingRepository defines the interface for parsing result storage
type ParsingRepository interface {
	SaveParsedWorkout(workout *domain.ParsedWorkout) error
	GetParsedWorkout(id string) (*domain.ParsedWorkout, error)
}

// EventPublisher defines the interface for publishing parsing events
type EventPublisher interface {
	PublishParsingSuccess(requestID string, workout *domain.ParsedWorkout) error
	PublishParsingFailure(requestID string, err error) error
}
