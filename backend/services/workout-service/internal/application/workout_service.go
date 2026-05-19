package application

import (
	"github.com/marcus/ironlog/backend/services/workout-service/internal/domain"
	"github.com/marcus/ironlog/backend/services/workout-service/internal/ports"
)

// WorkoutService handles workout business logic
type WorkoutService struct {
	workoutRepo ports.WorkoutRepository
	eventStore  ports.EventStore
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(repo ports.WorkoutRepository, eventStore ports.EventStore) *WorkoutService {
	return &WorkoutService{
		workoutRepo: repo,
		eventStore:  eventStore,
	}
}

// StartWorkoutCmd represents the command to start a workout
type StartWorkoutCmd struct {
	WorkoutID string
	UserID    string
	Notes     string
}

// StartWorkout initiates a new workout
func (ws *WorkoutService) StartWorkout(cmd StartWorkoutCmd) error {
	workout := domain.NewWorkout(cmd.WorkoutID, cmd.UserID)

	err := workout.StartWorkout(cmd.UserID, cmd.Notes)
	if err != nil {
		return err
	}

	// Persist events to event store
	err = ws.eventStore.SaveEvents(cmd.WorkoutID, workout.GetUncommittedEvents())
	if err != nil {
		return err
	}

	workout.MarkEventsAsCommitted()
	return nil
}

// PerformSetCmd represents the command to record a set
type PerformSetCmd struct {
	WorkoutID     string
	ExerciseID    string
	SetNumber     int
	SetType       string
	PlannedSets   int
	PlannedWeight float64
	ExecutedReps  []int
	TargetRepMin  int
	TargetRepMax  int
	Notes         string
}

// PerformSet records a completed set
func (ws *WorkoutService) PerformSet(cmd PerformSetCmd) error {
	// Load existing workout aggregate
	workout, err := ws.workoutRepo.GetByID(cmd.WorkoutID)
	if err != nil {
		return err
	}

	// Perform set within the aggregate
	err = workout.PerformSet(
		cmd.ExerciseID,
		cmd.SetNumber,
		cmd.SetType,
		cmd.PlannedSets,
		cmd.PlannedWeight,
		cmd.ExecutedReps,
		cmd.TargetRepMin,
		cmd.TargetRepMax,
		cmd.Notes,
	)
	if err != nil {
		return err
	}

	// Persist new events
	err = ws.eventStore.SaveEvents(cmd.WorkoutID, workout.GetUncommittedEvents())
	if err != nil {
		return err
	}

	workout.MarkEventsAsCommitted()
	return nil
}

// FinishWorkoutCmd represents the command to finish a workout
type FinishWorkoutCmd struct {
	WorkoutID string
	Notes     string
}

// FinishWorkout ends a workout session
func (ws *WorkoutService) FinishWorkout(cmd FinishWorkoutCmd) error {
	// Load existing workout aggregate
	workout, err := ws.workoutRepo.GetByID(cmd.WorkoutID)
	if err != nil {
		return err
	}

	// Finish the workout
	err = workout.FinishWorkout(cmd.Notes)
	if err != nil {
		return err
	}

	// Persist events
	err = ws.eventStore.SaveEvents(cmd.WorkoutID, workout.GetUncommittedEvents())
	if err != nil {
		return err
	}

	workout.MarkEventsAsCommitted()
	return nil
}
