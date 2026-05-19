package domain

// LoadSuggestion represents a suggested load increase
type LoadSuggestion struct {
	ExerciseID      string
	CurrentWeight   float64
	SuggestedWeight float64
	Reasoning       string
	ConfidenceScore float64
}

// DeloadSuggestion represents a suggested deload
type DeloadSuggestion struct {
	ExerciseID           string
	Reason               string
	RecommendedReduction float64
	Duration             string
}

// PerformanceRecommendation represents a general performance recommendation
type PerformanceRecommendation struct {
	ExerciseID  string
	Type        string // LOAD_INCREASE, DELOAD, FREQUENCY_INCREASE
	Description string
	Priority    string // HIGH, MEDIUM, LOW
}

// SuggestLoadIncrease determines if load should be increased
func SuggestLoadIncrease(repsAchieved []int, targetMin int, targetMax int) bool {
	if len(repsAchieved) == 0 {
		return false
	}

	// Count how many times the user hit the top of the range
	hits := 0
	for _, reps := range repsAchieved {
		if reps >= targetMax {
			hits++
		}
	}

	// If at least 2 sets hit the top, suggest increase
	return hits >= 2
}

// CalculateSuggestedIncrease calculates the suggested weight increase
func CalculateSuggestedIncrease(currentWeight float64) float64 {
	// Standard plate increment logic
	if currentWeight < 20 {
		return currentWeight + 2.5
	}
	if currentWeight < 60 {
		return currentWeight + 5
	}
	return currentWeight + 2.5 // Fine adjustments for heavier weights
}
