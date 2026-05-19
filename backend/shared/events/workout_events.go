package events

import "time"

// WorkoutStartedPayload represents the start of a workout
type WorkoutStartedPayload struct {
	WorkoutID string    `json:"workout_id"`
	UserID    string    `json:"user_id"`
	StartedAt time.Time `json:"started_at"`
	Notes     string    `json:"notes,omitempty"`
}

// ExerciseStartedPayload represents the start of an exercise within a workout
type ExerciseStartedPayload struct {
	WorkoutID  string    `json:"workout_id"`
	ExerciseID string    `json:"exercise_id"`
	Name       string    `json:"name"`
	StartedAt  time.Time `json:"started_at"`
}

// SetPerformedPayload represents a set completed in an exercise
type SetPerformedPayload struct {
	WorkoutID     string    `json:"workout_id"`
	ExerciseID    string    `json:"exercise_id"`
	SetNumber     int       `json:"set_number"`
	SetType       string    `json:"set_type"` // WARM_UP, FEEDER, WORK, TOP_SET
	PlannedSets   int       `json:"planned_sets"`
	PlannedWeight float64   `json:"planned_weight"`
	ExecutedReps  []int     `json:"executed_reps"`
	TargetRepMin  int       `json:"target_rep_min"`
	TargetRepMax  int       `json:"target_rep_max"`
	CompletedAt   time.Time `json:"completed_at"`
	Notes         string    `json:"notes,omitempty"`
}

// ExerciseFinishedPayload represents the end of an exercise
type ExerciseFinishedPayload struct {
	WorkoutID   string    `json:"workout_id"`
	ExerciseID  string    `json:"exercise_id"`
	TotalVolume float64   `json:"total_volume"` // weight * reps sum
	FinishedAt  time.Time `json:"finished_at"`
}

// WorkoutFinishedPayload represents the end of a workout
type WorkoutFinishedPayload struct {
	WorkoutID     string    `json:"workout_id"`
	UserID        string    `json:"user_id"`
	TotalVolume   float64   `json:"total_volume"`
	ExerciseCount int       `json:"exercise_count"`
	FinishedAt    time.Time `json:"finished_at"`
	Duration      int64     `json:"duration_seconds"`
	Notes         string    `json:"notes,omitempty"`
}

// WorkoutNoteAddedPayload represents a note added to a workout
type WorkoutNoteAddedPayload struct {
	WorkoutID string    `json:"workout_id"`
	Note      string    `json:"note"`
	AddedAt   time.Time `json:"added_at"`
}
