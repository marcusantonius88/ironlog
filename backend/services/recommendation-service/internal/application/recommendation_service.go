package application

import (
	"github.com/marcus/ironlog/backend/services/recommendation-service/internal/domain"
)

// RecommendationService handles recommendation logic
type RecommendationService struct{}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// GenerateLoadRecommendation generates a load increase recommendation
func (rs *RecommendationService) GenerateLoadRecommendation(exerciseID string, currentWeight float64, repsAchieved []int, targetMin, targetMax int) *domain.LoadSuggestion {
	if !domain.SuggestLoadIncrease(repsAchieved, targetMin, targetMax) {
		return nil
	}

	suggestedWeight := domain.CalculateSuggestedIncrease(currentWeight)

	return &domain.LoadSuggestion{
		ExerciseID:      exerciseID,
		CurrentWeight:   currentWeight,
		SuggestedWeight: suggestedWeight,
		Reasoning:       "Consistently hitting the top of the rep range",
		ConfidenceScore: 85.0,
	}
}

// GenerateDeloadRecommendation generates a deload recommendation
func (rs *RecommendationService) GenerateDeloadRecommendation(exerciseID string, regressionSessions int) *domain.DeloadSuggestion {
	if regressionSessions < 3 {
		return nil
	}

	return &domain.DeloadSuggestion{
		ExerciseID:           exerciseID,
		Reason:               "Performance regression detected over multiple sessions",
		RecommendedReduction: 10.0, // 10% reduction
		Duration:             "1 week",
	}
}
