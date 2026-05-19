package events

import "time"

// ParsingRequestedPayload represents a DSL parsing request
type ParsingRequestedPayload struct {
	RequestID   string    `json:"request_id"`
	UserID      string    `json:"user_id"`
	RawText     string    `json:"raw_text"`
	RequestedAt time.Time `json:"requested_at"`
}

// ParsedDSLPayload represents parsed DSL output
type ParsedDSLPayload struct {
	RequestID    string     `json:"request_id"`
	ExerciseName string     `json:"exercise_name"`
	SetGroups    []SetGroup `json:"set_groups"`
	RawText      string     `json:"raw_text"`
	ParsedAt     time.Time  `json:"parsed_at"`
}

// SetGroup represents a group of sets in the DSL
type SetGroup struct {
	SetType      string  `json:"set_type"` // WARM_UP, FEEDER, WORK, TOP_SET
	PlannedSets  int     `json:"planned_sets"`
	TargetRepMin int     `json:"target_rep_min"`
	TargetRepMax int     `json:"target_rep_max"`
	Weight       float64 `json:"weight"`
	ExecutedReps []int   `json:"executed_reps"`
}
