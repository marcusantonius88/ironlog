package domain

// ExerciseProgressionProjection represents read model for exercise progression
type ExerciseProgressionProjection struct {
	ExerciseID       string
	ExerciseName     string
	UserID           string
	CurrentLoad      float64
	LoadTrend        string
	RepsAccomplished int
	RepTrend         string
	VolumeTotal      float64
	VolumeTrend      string
	SessionsCount    int
	PersonalRecords  int
	LastPerformed    string // ISO timestamp
}

// WeeklyVolumeProjection represents weekly volume read model
type WeeklyVolumeProjection struct {
	UserID           string
	ExerciseID       string
	WeekStart        string // ISO date
	WeekEnd          string // ISO date
	TotalVolume      float64
	SessionCount     int
	AverageIntensity float64
}

// WorkoutTimelineProjection represents workout timeline read model
type WorkoutTimelineProjection struct {
	UserID        string
	WorkoutID     string
	Date          string // ISO date
	DurationMin   int
	TotalVolume   float64
	ExerciseCount int
	Exercises     []ExerciseSummary
	Notes         string
}

// ExerciseSummary represents brief exercise info in timeline
type ExerciseSummary struct {
	Name     string
	Volume   float64
	SetCount int
}

// PersonalRecordsProjection represents personal records read model
type PersonalRecordsProjection struct {
	UserID       string
	ExerciseID   string
	ExerciseName string
	RecordType   string // LOAD, REPS, VOLUME
	Value        float64
	AchievedAt   string // ISO timestamp
	Notes        string
}
