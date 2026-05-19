package http

import (
	"encoding/json"
	"net/http"

	"github.com/marcus/ironlog/backend/services/parser-service/internal/application"
)

// ParsingHandler handles HTTP requests for DSL parsing
type ParsingHandler struct {
	parsingService *application.ParsingService
}

// NewParsingHandler creates a new parsing handler
func NewParsingHandler(parsingService *application.ParsingService) *ParsingHandler {
	return &ParsingHandler{
		parsingService: parsingService,
	}
}

// ParseRequest represents the incoming DSL parsing request
type ParseRequest struct {
	RawText string `json:"raw_text"`
}

// ParseResponse represents the parsing response
type ParseResponse struct {
	Success   bool                `json:"success"`
	Exercise  string              `json:"exercise,omitempty"`
	SetGroups []ParsedSetGroupDTO `json:"set_groups,omitempty"`
	Error     string              `json:"error,omitempty"`
}

// ParsedSetGroupDTO represents a parsed set group in the response
type ParsedSetGroupDTO struct {
	SetType      string  `json:"set_type"`
	PlannedSets  int     `json:"planned_sets"`
	TargetRepMin int     `json:"target_rep_min"`
	TargetRepMax int     `json:"target_rep_max"`
	Weight       float64 `json:"weight"`
	Unit         string  `json:"unit"`
	ExecutedReps []int   `json:"executed_reps,omitempty"`
}

// ParseDSL handles the DSL parsing HTTP endpoint
func (ph *ParsingHandler) ParseDSL(w http.ResponseWriter, r *http.Request) {
	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse DSL
	workout, err := ph.parsingService.ParseDSL(req.RawText)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate
	if err := ph.parsingService.ValidateParsedWorkout(workout); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert to DTO
	setGroups := make([]ParsedSetGroupDTO, len(workout.SetGroups))
	for i, sg := range workout.SetGroups {
		setGroups[i] = ParsedSetGroupDTO{
			SetType:      string(sg.SetType),
			PlannedSets:  sg.PlannedSets,
			TargetRepMin: sg.TargetRepMin,
			TargetRepMax: sg.TargetRepMax,
			Weight:       sg.Weight,
			Unit:         sg.Unit,
			ExecutedReps: sg.ExecutedReps,
		}
	}

	response := ParseResponse{
		Success:   true,
		Exercise:  workout.ExerciseName,
		SetGroups: setGroups,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck returns the health status
func (ph *ParsingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
