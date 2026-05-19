package domain

import (
	"errors"
	"time"
)

// Workout represents the root aggregate for a training session
type Workout struct {
	ID            string
	UserID        string
	Status        WorkoutStatus
	ExerciseCount int
	TotalVolume   float64
	StartedAt     time.Time
	FinishedAt    *time.Time
	Duration      int64 // seconds
	Notes         string
	Version       int64
	Events        []interface{} // event sourcing
}

// WorkoutStatus enum
type WorkoutStatus string

const (
	WorkoutStatusActive    WorkoutStatus = "ACTIVE"
	WorkoutStatusFinished  WorkoutStatus = "FINISHED"
	WorkoutStatusCancelled WorkoutStatus = "CANCELLED"
)

// Exercise represents an exercise within a workout
type Exercise struct {
	ID          string
	WorkoutID   string
	Name        string
	TotalVolume float64
	SetCount    int
	StartedAt   time.Time
	FinishedAt  *time.Time
	Sets        []Set
	Version     int64
}

// Set represents a single set of an exercise
type Set struct {
	ID            string
	ExerciseID    string
	SetNumber     int
	SetType       string
	PlannedSets   int
	PlannedWeight float64
	ExecutedReps  []int
	TargetRepMin  int
	TargetRepMax  int
	CompletedAt   time.Time
	Notes         string
}

// Events

// WorkoutStartedEvent occurs when a workout session begins
type WorkoutStartedEvent struct {
	WorkoutID string
	UserID    string
	StartedAt time.Time
	Notes     string
}

// ExerciseStartedEvent occurs when an exercise begins
type ExerciseStartedEvent struct {
	ExerciseID string
	WorkoutID  string
	Name       string
	StartedAt  time.Time
}

// SetPerformedEvent occurs when a set is completed
type SetPerformedEvent struct {
	SetID         string
	ExerciseID    string
	WorkoutID     string
	SetNumber     int
	SetType       string
	PlannedSets   int
	PlannedWeight float64
	ExecutedReps  []int
	TargetRepMin  int
	TargetRepMax  int
	CompletedAt   time.Time
	Notes         string
}

// ExerciseFinishedEvent occurs when an exercise ends
type ExerciseFinishedEvent struct {
	ExerciseID  string
	WorkoutID   string
	TotalVolume float64
	FinishedAt  time.Time
}

// WorkoutFinishedEvent occurs when a workout session ends
type WorkoutFinishedEvent struct {
	WorkoutID     string
	UserID        string
	TotalVolume   float64
	ExerciseCount int
	FinishedAt    time.Time
	Duration      int64
	Notes         string
}

// WorkoutNoteAddedEvent occurs when a note is added to a workout
type WorkoutNoteAddedEvent struct {
	WorkoutID string
	Note      string
	AddedAt   time.Time
}

// Domain Errors
var (
	ErrWorkoutAlreadyStarted = errors.New("workout already started")
	ErrWorkoutNotActive      = errors.New("workout is not active")
	ErrInvalidSetData        = errors.New("invalid set data")
	ErrExerciseNotFound      = errors.New("exercise not found")
)

// NewWorkout creates a new workout aggregate
func NewWorkout(workoutID, userID string) *Workout {
	return &Workout{
		ID:        workoutID,
		UserID:    userID,
		Status:    WorkoutStatusActive,
		StartedAt: time.Now().UTC(),
		Events:    []interface{}{},
	}
}

// StartWorkout initializes a new workout
func (w *Workout) StartWorkout(userID string, notes string) error {
	if w.Status != "" {
		return ErrWorkoutAlreadyStarted
	}

	event := WorkoutStartedEvent{
		WorkoutID: w.ID,
		UserID:    userID,
		StartedAt: time.Now().UTC(),
		Notes:     notes,
	}

	w.applyEvent(event)
	w.Events = append(w.Events, event)
	return nil
}

// PerformSet records a completed set
func (w *Workout) PerformSet(exerciseID string, setNumber int, setType string, plannedSets int, plannedWeight float64, executedReps []int, targetRepMin, targetRepMax int, notes string) error {
	if w.Status != WorkoutStatusActive {
		return ErrWorkoutNotActive
	}

	setID := generateID("set")
	event := SetPerformedEvent{
		SetID:         setID,
		ExerciseID:    exerciseID,
		WorkoutID:     w.ID,
		SetNumber:     setNumber,
		SetType:       setType,
		PlannedSets:   plannedSets,
		PlannedWeight: plannedWeight,
		ExecutedReps:  executedReps,
		TargetRepMin:  targetRepMin,
		TargetRepMax:  targetRepMax,
		CompletedAt:   time.Now().UTC(),
		Notes:         notes,
	}

	w.applyEvent(event)
	w.Events = append(w.Events, event)
	return nil
}

// FinishWorkout ends the workout session
func (w *Workout) FinishWorkout(notes string) error {
	if w.Status != WorkoutStatusActive {
		return ErrWorkoutNotActive
	}

	now := time.Now().UTC()
	duration := int64(now.Sub(w.StartedAt).Seconds())

	event := WorkoutFinishedEvent{
		WorkoutID:     w.ID,
		UserID:        w.UserID,
		TotalVolume:   w.TotalVolume,
		ExerciseCount: w.ExerciseCount,
		FinishedAt:    now,
		Duration:      duration,
		Notes:         notes,
	}

	w.applyEvent(event)
	w.Events = append(w.Events, event)
	w.FinishedAt = &now
	w.Duration = duration
	return nil
}

// AddNote adds a note to the workout
func (w *Workout) AddNote(note string) error {
	event := WorkoutNoteAddedEvent{
		WorkoutID: w.ID,
		Note:      note,
		AddedAt:   time.Now().UTC(),
	}

	w.Events = append(w.Events, event)
	w.Notes += note + " "
	return nil
}

// applyEvent applies an event to the aggregate state
func (w *Workout) applyEvent(event interface{}) {
	switch e := event.(type) {
	case WorkoutStartedEvent:
		w.Status = WorkoutStatusActive
		w.StartedAt = e.StartedAt
		w.Notes = e.Notes
	case SetPerformedEvent:
		// Update exercise count if this is the first set
		if e.SetNumber == 1 {
			w.ExerciseCount++
		}
		// Calculate volume: weight * sum of reps
		repsSum := 0
		for _, rep := range e.ExecutedReps {
			repsSum += rep
		}
		volume := e.PlannedWeight * float64(repsSum)
		w.TotalVolume += volume
	case WorkoutFinishedEvent:
		w.Status = WorkoutStatusFinished
		w.FinishedAt = &e.FinishedAt
		w.Duration = e.Duration
	case WorkoutNoteAddedEvent:
		w.Notes += e.Note + " "
	}
	w.Version++
}

// GetUncommittedEvents returns events that have not been persisted
func (w *Workout) GetUncommittedEvents() []interface{} {
	return w.Events
}

// MarkEventsAsCommitted clears the uncommitted events
func (w *Workout) MarkEventsAsCommitted() {
	w.Events = []interface{}{}
}

// Helper function to generate unique IDs
func generateID(prefix string) string {
	// Implementation would use UUID or similar
	return prefix + "_" + time.Now().Format("20060102150405")
}
