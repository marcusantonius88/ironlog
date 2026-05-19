package domain

import (
	"time"
)

// PerformanceMetric tracks a metric over time
type PerformanceMetric struct {
	ExerciseID    string
	MetricType    string // LOAD, REPS, VOLUME
	CurrentValue  float64
	PreviousValue float64
	Trend         string // UP, DOWN, STABLE
	Sessions      int
	ChangedAt     time.Time
}

// PersonalRecord represents a personal record achievement
type PersonalRecord struct {
	ExerciseID string
	RecordType string // LOAD, REPS, VOLUME
	Value      float64
	AchievedAt time.Time
	SessionID  string
}

// DetectPerformanceChange detects improvement or regression
func DetectPerformanceChange(current, previous float64) string {
	if current > previous*1.02 { // 2% threshold
		return "UP"
	}
	if current < previous*0.98 {
		return "DOWN"
	}
	return "STABLE"
}

// CalculateVolume calculates total training volume
func CalculateVolume(weight float64, reps []int) float64 {
	repsSum := 0
	for _, r := range reps {
		repsSum += r
	}
	return weight * float64(repsSum)
}
