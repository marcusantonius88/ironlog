# Service Communication Guide - IRONLOG

## Overview

IRONLOG uses **event-driven architecture** where services communicate asynchronously through events published to Apache Kafka. All events follow a standardized `EventEnvelope` contract.

## Event Envelope Contract

Every event in IRONLOG has this structure:

```go
type EventEnvelope struct {
    EventID       string                `json:"event_id"`
    EventType     string                `json:"event_type"`
    AggregateID   string                `json:"aggregate_id"`
    CorrelationID string                `json:"correlation_id"`
    Payload       interface{}           `json:"payload"`
    CreatedAt     time.Time             `json:"created_at"`
    Metadata      map[string]string     `json:"metadata,omitempty"`
}
```

### Fields Explained

| Field | Purpose | Example |
|-------|---------|---------|
| **EventID** | Unique event identifier | `550e8400-e29b-41d4-a716-446655440000` |
| **EventType** | Type of event | `WorkoutStarted`, `SetPerformed` |
| **AggregateID** | Entity being changed (workout, exercise) | `wk_20240512_001` |
| **CorrelationID** | Links related events across services | `corr_user_session_123` |
| **Payload** | Event-specific data (varies by type) | JSON object |
| **CreatedAt** | ISO8601 timestamp | `2024-05-12T10:30:00Z` |
| **Metadata** | Optional key-value pairs | `{"user_id": "marcus"}` |

## Event Topics & Data Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          KAFKA TOPICS                                   │
└─────────────────────────────────────────────────────────────────────────┘

workout-events
    ├── WorkoutStarted          → projection-service, analytics-service
    ├── ExerciseStarted         → projection-service, analytics-service
    ├── SetPerformed            → analytics-service, projection-service
    └── WorkoutFinished         → analytics-service, projection-service

parser-events
    ├── DSLParsedSuccessfully   → none (frontend handles)
    └── DSLParsingFailed        → notification-service

analytics-events
    ├── PerformanceImproved     → recommendation-service, notification-service
    ├── PersonalRecordAchieved  → notification-service, projection-service
    └── TrendDetected           → recommendation-service

recommendation-events
    └── LoadIncreaseSuggested   → notification-service

notification-events
    └── NotificationSent        → none (end of chain)
```

## Service Responsibilities

### 1. **Workout Service** (Event Source)
**Ports:** 8080  
**Role:** Commands entry point. Receives workout commands and emits domain events.

**Publishes:**
- `WorkoutStarted`
- `ExerciseStarted`
- `SetPerformed`
- `WorkoutFinished`

**Example Event - SetPerformed:**
```json
{
  "event_id": "evt_001",
  "event_type": "SetPerformed",
  "aggregate_id": "wk_20240512_001",
  "correlation_id": "corr_user_session_123",
  "payload": {
    "workout_id": "wk_20240512_001",
    "exercise_id": "ex_supino_001",
    "set_number": 1,
    "set_type": "WORK",
    "planned_sets": 3,
    "planned_weight": 20.0,
    "executed_reps": [8, 7, 6],
    "target_rep_min": 6,
    "target_rep_max": 8,
    "unit": "kg"
  },
  "created_at": "2024-05-12T10:30:00Z"
}
```

### 2. **Parser Service** (Specialized Command Handler)
**Ports:** 8081  
**Role:** Parses DSL workout text into structured data.

**Consumes:** HTTP requests (not event-based)  
**Publishes:** Events to confirm/reject parsing

**Example Event - DSLParsedSuccessfully:**
```json
{
  "event_id": "evt_002",
  "event_type": "DSLParsedSuccessfully",
  "aggregate_id": "req_parse_001",
  "correlation_id": "corr_user_session_123",
  "payload": {
    "request_id": "req_parse_001",
    "exercise_name": "SUPINO RETO BARRA",
    "set_groups": [
      {
        "set_type": "WARM_UP",
        "planned_sets": 1,
        "target_rep_min": 1,
        "target_rep_max": 20,
        "weight": 10.0,
        "unit": "kg",
        "executed_reps": [10]
      }
    ]
  },
  "created_at": "2024-05-12T10:30:00Z"
}
```

### 3. **Analytics Service** (Event Processor)
**Ports:** 8082  
**Role:** Consumes workout events, computes metrics and insights.

**Consumes:**
- `SetPerformed` → Calculates volume, intensity metrics
- `WorkoutFinished` → Calculates totals, duration
- `ExerciseStarted` → Tracks timing

**Publishes:**
- `PerformanceImproved` (when metrics improve)
- `PersonalRecordAchieved` (new max weight/reps)
- `TrendDetected` (pattern recognized)

**Example Event - PersonalRecordAchieved:**
```json
{
  "event_id": "evt_003",
  "event_type": "PersonalRecordAchieved",
  "aggregate_id": "wk_20240512_001",
  "correlation_id": "corr_user_session_123",
  "payload": {
    "user_id": "marcus",
    "exercise_name": "SUPINO RETO",
    "previous_record": 25.0,
    "new_record": 27.5,
    "unit": "kg",
    "date_achieved": "2024-05-12",
    "notes": "Top set at 27.5kg x 3 reps"
  },
  "created_at": "2024-05-12T10:45:00Z"
}
```

### 4. **Projection Service** (CQRS Read Model Builder)
**Ports:** 8083  
**Role:** Builds and maintains read models for queries (CQRS pattern).

**Consumes:** All workout events  
**Publishes:** None (reads only)

**Read Models Maintained:**
- `exercise_progression` - Current weight/reps for each exercise
- `weekly_volume` - Total volume per week
- `workout_timeline` - Historical workout data
- `personal_records` - PRs for each exercise

**Query Example:**
```sql
-- Get current max weight for exercise
SELECT MAX(planned_weight) 
FROM exercise_progression 
WHERE exercise_id = 'ex_supino_001' 
AND date >= CURRENT_DATE - INTERVAL '30 days';
```

### 5. **Recommendation Service** (Insights Engine)
**Ports:** 8084  
**Role:** Analyzes trends and suggests progression.

**Consumes:**
- `PerformanceImproved` → Detects when to increase load
- `TrendDetected` → Pattern analysis

**Publishes:**
- `LoadIncreaseSuggested`

**Example Event - LoadIncreaseSuggested:**
```json
{
  "event_id": "evt_004",
  "event_type": "LoadIncreaseSuggested",
  "aggregate_id": "ex_supino_001",
  "correlation_id": "corr_user_session_123",
  "payload": {
    "user_id": "marcus",
    "exercise_id": "ex_supino_001",
    "exercise_name": "SUPINO RETO",
    "current_weight": 25.0,
    "suggested_weight": 27.5,
    "reason": "Successfully completed 3x8 for 3 consecutive sessions",
    "confidence": 0.95
  },
  "created_at": "2024-05-12T11:00:00Z"
}
```

### 6. **Notification Service** (Alerting)
**Ports:** 8085  
**Role:** Consumes interesting events and dispatches notifications.

**Consumes:**
- `PersonalRecordAchieved`
- `LoadIncreaseSuggested`
- `DSLParsingFailed`

**Publishes:**
- `NotificationSent` (final event)

**Example Event - NotificationSent:**
```json
{
  "event_id": "evt_005",
  "event_type": "NotificationSent",
  "aggregate_id": "notif_001",
  "correlation_id": "corr_user_session_123",
  "payload": {
    "notification_id": "notif_001",
    "user_id": "marcus",
    "type": "PR_ACHIEVEMENT",
    "title": "New Personal Record!",
    "message": "You achieved a new PR on SUPINO RETO: 27.5kg x 3!",
    "sent_at": "2024-05-12T11:00:01Z",
    "delivery_status": "sent"
  },
  "created_at": "2024-05-12T11:00:01Z"
}
```

## Communication Patterns

### Pattern 1: Request-Response (Frontend → Services)

```
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │ HTTP POST /parse
       │ {"raw_text": "..."}
       ▼
┌──────────────────┐
│  Parser Service  │
└──────┬───────────┘
       │ HTTP 200
       │ {"success": true, "exercise_name": "..."}
       ▼
┌─────────────┐
│   Frontend  │
└─────────────┘
```

### Pattern 2: Event-Driven (Async)

```
┌──────────────────┐
│ Workout Service  │ → Publishes WorkoutStarted
└──────────────────┘
         ▼
     ┌────────┐
     │ Kafka  │
     └────────┘
      ▼      ▼      ▼
┌──────────┐ ┌──────────────┐ ┌─────────────┐
│Projection│ │  Analytics   │ │Notification │
│ Service  │ │   Service    │ │  Service    │
└──────────┘ └──────────────┘ └─────────────┘
```

### Pattern 3: Event-Triggered Pipeline

```
SetPerformed Event
    ▼
Analytics Service
    ├─→ Calculates metrics
    ├─→ Publishes PerformanceImproved
    ▼
Kafka (analytics-events topic)
    ▼
Recommendation Service
    ├─→ Analyzes trends
    ├─→ Publishes LoadIncreaseSuggested
    ▼
Kafka (recommendation-events topic)
    ▼
Notification Service
    ├─→ Prepares message
    ├─→ Publishes NotificationSent
    ▼
Database (event store)
```

## Kafka Topics Configuration

| Topic | Partitions | Purpose |
|-------|-----------|---------|
| `workout-events` | 3 | Core domain events from workouts |
| `parser-events` | 1 | DSL parsing results |
| `analytics-events` | 2 | Computed insights |
| `recommendation-events` | 1 | Suggested improvements |
| `notification-events` | 1 | Notification dispatch |

## Idempotency Strategy

Services use **Redis** to track processed events and prevent duplicates:

```
Event arrives → Check Redis for EventID
                ├→ Found: Skip (already processed)
                └→ Not found: Process + Store in Redis with TTL
```

```bash
# Check for processed event
REDIS: GET "event:evt_001"  
# Result: "processed" or nil

# Store processed event
REDIS: SET "event:evt_001" "processed" EX 86400  # 24h TTL
```

## Database Storage

All events are persisted in PostgreSQL `event_store` table:

```sql
CREATE TABLE event_store (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR NOT NULL,
    aggregate_id VARCHAR NOT NULL,
    correlation_id VARCHAR,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    metadata JSONB,
    INDEX (aggregate_id),
    INDEX (event_type),
    INDEX (created_at)
);
```

## Failure Handling

### Dead Letter Queue (DLQ)

Events that fail processing are sent to DLQ topic:

```
Topic: workout-events-dlq
Reason: max_retries_exceeded | deserialization_error | processing_timeout
```

```bash
# Monitor DLQ
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic workout-events-dlq \
  --from-beginning
```

### Retry Strategy

- **Immediate retry**: First 3 attempts (on same service)
- **Exponential backoff**: Retry after 5s, 10s, 30s
- **Dead letter**: After 5 total attempts, move to DLQ

## Monitoring Communication

### View Event Flow

```bash
# Monitor all workout events in real-time
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic workout-events \
  --from-beginning \
  --property print.timestamp=true

# Monitor specific consumer group lag
docker-compose exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --group analytics-service-group \
  --describe
```

### Query Event Store

```bash
# Recent events
SELECT event_id, event_type, aggregate_id, created_at 
FROM event_store 
ORDER BY created_at DESC 
LIMIT 100;

# Events for specific workout
SELECT event_type, payload, created_at 
FROM event_store 
WHERE aggregate_id = 'wk_20240512_001' 
ORDER BY created_at;

# Count by event type
SELECT event_type, COUNT(*) 
FROM event_store 
GROUP BY event_type 
ORDER BY COUNT(*) DESC;
```

## Sequence Diagram: Complete Workout Flow

```
User          Frontend      Parser      Workout      Analytics   Notification
│                │             │            │            │            │
├─ Submit DSL ──>│             │            │            │            │
│                │─ POST /parse─>           │            │            │
│                │             │ Parse DSL  │            │            │
│                │<─ Parsed ────│            │            │            │
│                │ Exercise ──>│            │            │            │
│                │              │─ Record ─>│            │            │
│                │              │  Workout  │            │            │
│                │              │           ├─ Start ─>──┐            │
│                │              │           │           (store)       │
│                │              │           │            │            │
│                ├─ Start Set ─────────────>│            │            │
│                │              │           │            │            │
│                ├─ Finish Set ─────────────>│            │            │
│                │              │           │─ Set ─────>│            │
│                │              │           │ Performed  │            │
│                │              │           │            ├─ Analyze ─>│
│                │              │           │            │ Update     │
│                │              │           │            │ PRs        │
│                │              │           │            ├─ PR ─────>│
│                │              │           │            │ Achieved  │
│                │              │           │            │           │
│                │              │           │            │    ├─ Send ──>
│                │              │           │            │    │
│<─ Notification ───────────────────────────────────────────────┤
│                │             │            │            │            │
```

## Next Steps

- Review [ARCHITECTURE.md](ARCHITECTURE.md) for design patterns
- Check [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) to add new events
- See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues
