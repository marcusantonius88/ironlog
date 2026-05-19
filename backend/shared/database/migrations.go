package database

// Migration SQL scripts for event sourcing infrastructure

const EventStoreTableSQL = `
CREATE TABLE IF NOT EXISTS event_store (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    correlation_id UUID NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    INDEX idx_aggregate_id (aggregate_id),
    INDEX idx_created_at (created_at),
    INDEX idx_event_type (event_type)
);
`

const OutboxTableSQL = `
CREATE TABLE IF NOT EXISTS outbox (
    id SERIAL PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL,
    event_id UUID NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    correlation_id UUID NOT NULL,
    payload JSONB NOT NULL,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP,
    INDEX idx_published (published),
    INDEX idx_created_at (created_at)
);
`

const SnapshotTableSQL = `
CREATE TABLE IF NOT EXISTS snapshots (
    snapshot_id UUID PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL UNIQUE,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_version BIGINT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    INDEX idx_aggregate_id (aggregate_id)
);
`

const IdempotencyTableSQL = `
CREATE TABLE IF NOT EXISTS idempotency_keys (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    request_id UUID NOT NULL,
    response JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    INDEX idx_expires_at (expires_at)
);
`

const WorkoutAggregateTableSQL = `
CREATE TABLE IF NOT EXISTS workout_aggregate (
    workout_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_volume DECIMAL(10,2),
    exercise_count INT,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    duration_seconds BIGINT,
    notes TEXT,
    version BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
);
`

const ExerciseAggregateTableSQL = `
CREATE TABLE IF NOT EXISTS exercise_aggregate (
    exercise_id VARCHAR(255) PRIMARY KEY,
    workout_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    total_volume DECIMAL(10,2),
    set_count INT,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    version BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (workout_id) REFERENCES workout_aggregate(workout_id) ON DELETE CASCADE,
    INDEX idx_workout_id (workout_id)
);
`

const ProjectionExerciseProgressionTableSQL = `
CREATE TABLE IF NOT EXISTS exercise_progression_projection (
    id SERIAL PRIMARY KEY,
    exercise_id VARCHAR(255) NOT NULL,
    exercise_name VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    current_load DECIMAL(10,2),
    load_trend VARCHAR(50),
    reps_accomplished INT,
    rep_trend VARCHAR(50),
    volume_total DECIMAL(12,2),
    volume_trend VARCHAR(50),
    sessions_count INT,
    personal_records INT,
    last_performed TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(exercise_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_exercise_name (exercise_name)
);
`

const ProjectionWeeklyVolumeTableSQL = `
CREATE TABLE IF NOT EXISTS weekly_volume_projection (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    exercise_id VARCHAR(255),
    week_start DATE NOT NULL,
    week_end DATE NOT NULL,
    total_volume DECIMAL(12,2),
    session_count INT,
    average_intensity DECIMAL(5,2),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, exercise_id, week_start),
    INDEX idx_user_id (user_id),
    INDEX idx_week_start (week_start)
);
`

const ProjectionWorkoutTimelineTableSQL = `
CREATE TABLE IF NOT EXISTS workout_timeline_projection (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    workout_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    duration_minutes INT,
    total_volume DECIMAL(12,2),
    exercise_count INT,
    exercises JSONB,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(workout_id),
    INDEX idx_user_id (user_id),
    INDEX idx_date (date)
);
`

const ProjectionPersonalRecordsTableSQL = `
CREATE TABLE IF NOT EXISTS personal_records_projection (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    exercise_id VARCHAR(255) NOT NULL,
    exercise_name VARCHAR(255) NOT NULL,
    record_type VARCHAR(50) NOT NULL,
    value DECIMAL(12,2) NOT NULL,
    achieved_at TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    INDEX idx_user_id (user_id),
    INDEX idx_exercise_id (exercise_id),
    INDEX idx_achieved_at (achieved_at)
);
`
