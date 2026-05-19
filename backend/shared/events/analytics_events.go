package events

import "time"

// PerformanceImprovedPayload indicates performance improvement
type PerformanceImprovedPayload struct {
	AggregateID      string    `json:"aggregate_id"` // exercise_id or workout_id
	ExerciseName     string    `json:"exercise_name"`
	MetricType       string    `json:"metric_type"` // LOAD, REPS, VOLUME
	PreviousValue    float64   `json:"previous_value"`
	CurrentValue     float64   `json:"current_value"`
	ImpromentPercent float64   `json:"improvement_percent"`
	ConsecutiveSets  int       `json:"consecutive_sets"`
	DetectedAt       time.Time `json:"detected_at"`
}

// PerformanceRegressedPayload indicates performance regression
type PerformanceRegressedPayload struct {
	AggregateID       string    `json:"aggregate_id"`
	ExerciseName      string    `json:"exercise_name"`
	MetricType        string    `json:"metric_type"`
	PreviousValue     float64   `json:"previous_value"`
	CurrentValue      float64   `json:"current_value"`
	RegressionPercent float64   `json:"regression_percent"`
	Sessions          int       `json:"sessions"`
	DetectedAt        time.Time `json:"detected_at"`
}

// PersonalRecordAchievedPayload indicates a new personal record
type PersonalRecordAchievedPayload struct {
	AggregateID    string    `json:"aggregate_id"`
	ExerciseName   string    `json:"exercise_name"`
	RecordType     string    `json:"record_type"` // LOAD, REPS, VOLUME
	Value          float64   `json:"value"`
	PreviousRecord float64   `json:"previous_record"`
	AchievedAt     time.Time `json:"achieved_at"`
}

// VolumeIncreasedPayload indicates volume progression
type VolumeIncreasedPayload struct {
	AggregateID     string    `json:"aggregate_id"`
	ExerciseName    string    `json:"exercise_name"`
	PreviousVolume  float64   `json:"previous_volume"`
	CurrentVolume   float64   `json:"current_volume"`
	IncreasePercent float64   `json:"increase_percent"`
	Period          string    `json:"period"` // WEEKLY, MONTHLY, ALL_TIME
	DetectedAt      time.Time `json:"detected_at"`
}

// LoadIncreaseSuggestedPayload suggests load increase
type LoadIncreaseSuggestedPayload struct {
	AggregateID     string    `json:"aggregate_id"`
	ExerciseName    string    `json:"exercise_name"`
	CurrentWeight   float64   `json:"current_weight"`
	SuggestedWeight float64   `json:"suggested_weight"`
	Reasoning       string    `json:"reasoning"`
	ConfidenceScore float64   `json:"confidence_score"` // 0-100
	SuggestedAt     time.Time `json:"suggested_at"`
}
