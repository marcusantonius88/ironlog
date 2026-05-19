-- Create event store table
CREATE TABLE IF NOT EXISTS event_store (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    correlation_id UUID NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_event_store_aggregate_id ON event_store(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_event_store_created_at ON event_store(created_at);
CREATE INDEX IF NOT EXISTS idx_event_store_event_type ON event_store(event_type);

-- Create outbox table for Outbox Pattern + CDC
CREATE TABLE IF NOT EXISTS outbox (
    id SERIAL PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL,
    event_id UUID NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    correlation_id UUID NOT NULL,
    payload JSONB NOT NULL,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_published ON outbox(published);
CREATE INDEX IF NOT EXISTS idx_outbox_created_at ON outbox(created_at);

-- Create snapshots table for performance optimization
CREATE TABLE IF NOT EXISTS snapshots (
    snapshot_id UUID PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL UNIQUE,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_version BIGINT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_snapshots_aggregate_id ON snapshots(aggregate_id);

-- Create idempotency keys table
CREATE TABLE IF NOT EXISTS idempotency_keys (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    request_id UUID NOT NULL,
    response JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_idempotency_expires_at ON idempotency_keys(expires_at);

-- Create workout aggregate table
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
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workout_aggregate_user_id ON workout_aggregate(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_aggregate_status ON workout_aggregate(status);

-- Create exercise aggregate table
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
    FOREIGN KEY (workout_id) REFERENCES workout_aggregate(workout_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_exercise_aggregate_workout_id ON exercise_aggregate(workout_id);

-- Create CQRS projection tables
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
    UNIQUE(exercise_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_exercise_progression_user_id ON exercise_progression_projection(user_id);
CREATE INDEX IF NOT EXISTS idx_exercise_progression_exercise_name ON exercise_progression_projection(exercise_name);

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
    UNIQUE(user_id, exercise_id, week_start)
);

CREATE INDEX IF NOT EXISTS idx_weekly_volume_user_id ON weekly_volume_projection(user_id);
CREATE INDEX IF NOT EXISTS idx_weekly_volume_week_start ON weekly_volume_projection(week_start);

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
    UNIQUE(workout_id)
);

CREATE INDEX IF NOT EXISTS idx_workout_timeline_user_id ON workout_timeline_projection(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_timeline_date ON workout_timeline_projection(date);

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
    UNIQUE(user_id, exercise_id, record_type)
);

CREATE INDEX IF NOT EXISTS idx_personal_records_user_id ON personal_records_projection(user_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_exercise_id ON personal_records_projection(exercise_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_achieved_at ON personal_records_projection(achieved_at);
