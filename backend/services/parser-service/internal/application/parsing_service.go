package application

import (
	"github.com/marcus/ironlog/backend/services/parser-service/internal/domain"
)

// ParsingService orchestrates the parsing workflow
type ParsingService struct{}

// NewParsingService creates a new parsing service
func NewParsingService() *ParsingService {
	return &ParsingService{}
}

// ParseDSL parses raw DSL text and returns structured data
func (ps *ParsingService) ParseDSL(rawText string) (*domain.ParsedWorkout, error) {
	// Tokenize
	lexer := domain.NewLexer(rawText)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	// Parse
	parser := domain.NewParser(tokens)
	workout, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	workout.RawInput = rawText

	return workout, nil
}

// ValidateParsedWorkout validates parsed workout structure
func (ps *ParsingService) ValidateParsedWorkout(workout *domain.ParsedWorkout) error {
	if workout.ExerciseName == "" {
		return domain.ParseError{Message: "exercise name cannot be empty", Line: 1, Col: 1}
	}

	if len(workout.SetGroups) == 0 {
		return domain.ParseError{Message: "at least one set group is required", Line: 2, Col: 1}
	}

	for i, sg := range workout.SetGroups {
		if sg.PlannedSets <= 0 {
			return domain.ParseError{Message: "planned sets must be greater than 0", Line: i + 2, Col: 1}
		}

		if sg.TargetRepMin <= 0 || sg.TargetRepMax <= 0 {
			return domain.ParseError{Message: "rep ranges must be positive", Line: i + 2, Col: 1}
		}

		if sg.TargetRepMin > sg.TargetRepMax {
			return domain.ParseError{Message: "min reps cannot exceed max reps", Line: i + 2, Col: 1}
		}

		if sg.Weight <= 0 {
			return domain.ParseError{Message: "weight must be positive", Line: i + 2, Col: 1}
		}
	}

	return nil
}
