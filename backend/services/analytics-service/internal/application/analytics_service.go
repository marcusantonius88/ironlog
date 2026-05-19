package application

import (
	"github.com/marcus/ironlog/backend/services/analytics-service/internal/domain"
)

// AnalyticsService handles analytics computation
type AnalyticsService struct{}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// DetectPersonalRecord detects if a new PR was achieved
func (as *AnalyticsService) DetectPersonalRecord(exerciseID string, recordType string, currentValue float64, previousRecord float64) bool {
	return currentValue > previousRecord
}

// AnalyzePerformanceTrend analyzes performance trend
func (as *AnalyticsService) AnalyzePerformanceTrend(exerciseID string, values []float64) string {
	if len(values) < 2 {
		return "INSUFFICIENT_DATA"
	}

	improvements := 0
	regressions := 0

	for i := 1; i < len(values); i++ {
		trend := domain.DetectPerformanceChange(values[i], values[i-1])
		if trend == "UP" {
			improvements++
		} else if trend == "DOWN" {
			regressions++
		}
	}

	if improvements > regressions {
		return "IMPROVING"
	} else if regressions > improvements {
		return "REGRESSING"
	}
	return "STABLE"
}

// CalculateWeeklyVolume calculates total volume for a week
func (as *AnalyticsService) CalculateWeeklyVolume(sets []SetData) float64 {
	totalVolume := 0.0
	for _, set := range sets {
		totalVolume += domain.CalculateVolume(set.Weight, set.ExecutedReps)
	}
	return totalVolume
}

// SetData represents set information for calculation
type SetData struct {
	Weight       float64
	ExecutedReps []int
}
