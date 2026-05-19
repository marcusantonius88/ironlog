package kafka

import (
	"encoding/json"
	"log"

	"github.com/marcus/ironlog/backend/services/parser-service/internal/domain"
	"github.com/marcus/ironlog/backend/shared/events"
	"github.com/marcus/ironlog/backend/shared/infra"
)

// EventPublisher publishes parsing-related events to Kafka
type EventPublisher struct {
	kafkaProducer *infra.KafkaProducer
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(producer *infra.KafkaProducer) *EventPublisher {
	return &EventPublisher{
		kafkaProducer: producer,
	}
}

// PublishParsingSuccess publishes a successful parsing event
func (ep *EventPublisher) PublishParsingSuccess(requestID string, workout *domain.ParsedWorkout) error {
	setGroups := make([]events.SetGroup, len(workout.SetGroups))
	for i, sg := range workout.SetGroups {
		setGroups[i] = events.SetGroup{
			SetType:      string(sg.SetType),
			PlannedSets:  sg.PlannedSets,
			TargetRepMin: sg.TargetRepMin,
			TargetRepMax: sg.TargetRepMax,
			Weight:       sg.Weight,
			ExecutedReps: sg.ExecutedReps,
		}
	}

	payload := events.ParsedDSLPayload{
		RequestID:    requestID,
		ExerciseName: workout.ExerciseName,
		SetGroups:    setGroups,
		RawText:      workout.RawInput,
	}

	event := events.NewEventEnvelope("DSLParsedSuccessfully", requestID, requestID, payload)

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return err
	}

	return ep.kafkaProducer.PublishEvent(requestID, event)
}

// PublishParsingFailure publishes a parsing failure event
func (ep *EventPublisher) PublishParsingFailure(requestID string, err error) error {
	failurePayload := map[string]interface{}{
		"request_id": requestID,
		"error":      err.Error(),
	}

	event := events.NewEventEnvelope("DSLParsingFailed", requestID, requestID, failurePayload)

	return ep.kafkaProducer.PublishEvent(requestID, event)
}
