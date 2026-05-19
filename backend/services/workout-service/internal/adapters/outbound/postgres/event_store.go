package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/marcus/ironlog/backend/services/workout-service/internal/domain"
)

// PostgresEventStore implements event sourcing with PostgreSQL
type PostgresEventStore struct {
	db *sql.DB
}

// NewPostgresEventStore creates a new PostgreSQL event store
func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

// SaveEvents persists events to the event store and outbox (Outbox Pattern)
func (pes *PostgresEventStore) SaveEvents(aggregateID string, events []interface{}) error {
	tx, err := pes.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		eventID := uuid.New().String()
		eventType := getEventType(event)
		correlationID := uuid.New().String()

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		// Save to event store
		_, err = tx.Exec(`
			INSERT INTO event_store (event_id, event_type, aggregate_id, aggregate_type, correlation_id, payload, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, eventID, eventType, aggregateID, "Workout", correlationID, payload, time.Now().UTC())
		if err != nil {
			log.Printf("failed to save event: %v", err)
			return err
		}

		// Save to outbox for eventual publishing (Outbox Pattern)
		_, err = tx.Exec(`
			INSERT INTO outbox (aggregate_id, event_id, event_type, correlation_id, payload, published, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, aggregateID, eventID, eventType, correlationID, payload, false, time.Now().UTC())
		if err != nil {
			log.Printf("failed to save outbox event: %v", err)
			return err
		}
	}

	return tx.Commit()
}

// LoadEvents retrieves all events for an aggregate
func (pes *PostgresEventStore) LoadEvents(aggregateID string) ([]interface{}, error) {
	rows, err := pes.db.Query(`
		SELECT payload, event_type FROM event_store
		WHERE aggregate_id = $1
		ORDER BY created_at ASC
	`, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []interface{}
	for rows.Next() {
		var payload []byte
		var eventType string
		if err := rows.Scan(&payload, &eventType); err != nil {
			return nil, err
		}

		// Deserialize to appropriate event type
		event, err := deserializeEvent(eventType, payload)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventStream returns the event stream for an aggregate
func (pes *PostgresEventStore) GetEventStream(aggregateID string) ([]interface{}, error) {
	return pes.LoadEvents(aggregateID)
}

func getEventType(event interface{}) string {
	switch event.(type) {
	case domain.WorkoutStartedEvent:
		return "WorkoutStarted"
	case domain.SetPerformedEvent:
		return "SetPerformed"
	case domain.ExerciseStartedEvent:
		return "ExerciseStarted"
	case domain.ExerciseFinishedEvent:
		return "ExerciseFinished"
	case domain.WorkoutFinishedEvent:
		return "WorkoutFinished"
	case domain.WorkoutNoteAddedEvent:
		return "WorkoutNoteAdded"
	default:
		return "Unknown"
	}
}

func deserializeEvent(eventType string, payload []byte) (interface{}, error) {
	switch eventType {
	case "WorkoutStarted":
		var event domain.WorkoutStartedEvent
		json.Unmarshal(payload, &event)
		return event, nil
	case "SetPerformed":
		var event domain.SetPerformedEvent
		json.Unmarshal(payload, &event)
		return event, nil
	case "WorkoutFinished":
		var event domain.WorkoutFinishedEvent
		json.Unmarshal(payload, &event)
		return event, nil
	default:
		var event map[string]interface{}
		json.Unmarshal(payload, &event)
		return event, nil
	}
}
