package http

import (
	"encoding/json"
	"net/http"

	"github.com/marcus/ironlog/backend/services/workout-service/internal/application"
)

// WorkoutHandler handles HTTP requests for workout operations
type WorkoutHandler struct {
	workoutService *application.WorkoutService
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(service *application.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: service,
	}
}

// StartWorkoutRequest represents a start workout request
type StartWorkoutRequest struct {
	WorkoutID string `json:"workout_id"`
	UserID    string `json:"user_id"`
	Notes     string `json:"notes,omitempty"`
}

// StartWorkout handles the start workout endpoint
func (wh *WorkoutHandler) StartWorkout(w http.ResponseWriter, r *http.Request) {
	var req StartWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	cmd := application.StartWorkoutCmd{
		WorkoutID: req.WorkoutID,
		UserID:    req.UserID,
		Notes:     req.Notes,
	}

	if err := wh.workoutService.StartWorkout(cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "workout started"})
}

// PerformSetRequest represents a set performance request
type PerformSetRequest struct {
	WorkoutID     string  `json:"workout_id"`
	ExerciseID    string  `json:"exercise_id"`
	SetNumber     int     `json:"set_number"`
	SetType       string  `json:"set_type"`
	PlannedSets   int     `json:"planned_sets"`
	PlannedWeight float64 `json:"planned_weight"`
	ExecutedReps  []int   `json:"executed_reps"`
	TargetRepMin  int     `json:"target_rep_min"`
	TargetRepMax  int     `json:"target_rep_max"`
	Notes         string  `json:"notes,omitempty"`
}

// PerformSet handles the perform set endpoint
func (wh *WorkoutHandler) PerformSet(w http.ResponseWriter, r *http.Request) {
	var req PerformSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	cmd := application.PerformSetCmd{
		WorkoutID:     req.WorkoutID,
		ExerciseID:    req.ExerciseID,
		SetNumber:     req.SetNumber,
		SetType:       req.SetType,
		PlannedSets:   req.PlannedSets,
		PlannedWeight: req.PlannedWeight,
		ExecutedReps:  req.ExecutedReps,
		TargetRepMin:  req.TargetRepMin,
		TargetRepMax:  req.TargetRepMax,
		Notes:         req.Notes,
	}

	if err := wh.workoutService.PerformSet(cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "set recorded"})
}

// FinishWorkoutRequest represents a finish workout request
type FinishWorkoutRequest struct {
	WorkoutID string `json:"workout_id"`
	Notes     string `json:"notes,omitempty"`
}

// FinishWorkout handles the finish workout endpoint
func (wh *WorkoutHandler) FinishWorkout(w http.ResponseWriter, r *http.Request) {
	var req FinishWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	cmd := application.FinishWorkoutCmd{
		WorkoutID: req.WorkoutID,
		Notes:     req.Notes,
	}

	if err := wh.workoutService.FinishWorkout(cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "workout finished"})
}

// HealthCheck returns the health status
func (wh *WorkoutHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
