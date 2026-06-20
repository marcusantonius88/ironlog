# Development Guide - IRONLOG

## Quick Reference

This guide explains how to develop new features in IRONLOG. We'll walk through a complete example: **Adding a new workout metric**.

## Prerequisites

- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- Familiarity with [ARCHITECTURE.md](ARCHITECTURE.md)
- Understanding of event-driven systems

## Development Environment Setup

### 1. Start Infrastructure

```bash
cd infra
docker-compose up -d
```

Wait for all services to be ready:

```bash
docker-compose ps  # All containers should be "Up"
```

### 2. Test Infrastructure

```bash
# Check database
docker-compose exec postgres psql -U postgres -d ironlog -c "\dt"

# Check Kafka
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Check services
curl http://localhost:8081/health  # Parser
curl http://localhost:8080/health  # Workout
```

## Development Workflow

### Step 1: Define Your Event

All new functionality starts with an event. Events are the facts that happened.

**File:** `backend/shared/events/analytics_events.go`

```go
// New event for a metric we want to track
type StrengthProgressionPayload struct {
    ExerciseID          string    `json:"exercise_id"`
    ExerciseName        string    `json:"exercise_name"`
    PreviousBestWeight  float64   `json:"previous_best_weight"`
    NewBestWeight       float64   `json:"new_best_weight"`
    Unit                string    `json:"unit"`
    ConsecutiveSessions int       `json:"consecutive_sessions"`
    DetectedAt          time.Time `json:"detected_at"`
}
```

**Why start with events?**
- Defines what data flows through the system
- Used by all services that need this information
- Acts as a contract between services
- Makes testing easier

### Step 2: Implement Domain Logic

Domain logic lives in the `domain/` folder and should be **framework-agnostic**.

**File:** `backend/services/analytics-service/internal/domain/strength_progression.go`

```go
package domain

import "errors"

type StrengthProgression struct {
    ExerciseID              string
    ExerciseName            string
    PreviousBestWeight      float64
    NewBestWeight           float64
    Unit                    string
    ConsecutiveSessionCount int
}

// Calculate if this is a significant progression
func (sp *StrengthProgression) IsSignificantProgression() bool {
    // At least 5% improvement
    threshold := sp.PreviousBestWeight * 0.05
    improvement := sp.NewBestWeight - sp.PreviousBestWeight
    
    return improvement >= threshold && sp.ConsecutiveSessionCount >= 3
}

// Validate data
func (sp *StrengthProgression) Validate() error {
    if sp.ExerciseID == "" {
        return errors.New("exercise_id is required")
    }
    if sp.NewBestWeight <= 0 {
        return errors.New("new_best_weight must be positive")
    }
    if sp.NewBestWeight < sp.PreviousBestWeight {
        return errors.New("new weight cannot be less than previous")
    }
    return nil
}
```

**Key principles:**
- Pure functions (same input = same output)
- No external dependencies (no Kafka, DB, HTTP)
- Test-friendly (no mocking needed)
- Business logic only

### Step 3: Create Application Service

The application service orchestrates the domain logic and coordinates with ports.

**File:** `backend/services/analytics-service/internal/application/strength_progression_service.go`

```go
package application

import (
    "context"
    "github.com/google/uuid"
    "github.com/marcus/ironlog/backend/services/analytics-service/internal/domain"
    "github.com/marcus/ironlog/backend/services/analytics-service/internal/ports"
    "github.com/marcus/ironlog/backend/shared/events"
    "time"
)

type StrengthProgressionService struct {
    eventPublisher ports.EventPublisher
    repository     ports.Repository
}

func NewStrengthProgressionService(
    eventPublisher ports.EventPublisher,
    repository ports.Repository,
) *StrengthProgressionService {
    return &StrengthProgressionService{
        eventPublisher: eventPublisher,
        repository:     repository,
    }
}

// DetectStrengthProgression checks if there's significant strength gain
func (s *StrengthProgressionService) DetectStrengthProgression(
    ctx context.Context,
    exerciseID string,
    exerciseName string,
    previousBest float64,
    currentBest float64,
    unit string,
    sessions int,
) error {
    // Create domain entity
    progression := &domain.StrengthProgression{
        ExerciseID:              exerciseID,
        ExerciseName:            exerciseName,
        PreviousBestWeight:      previousBest,
        NewBestWeight:           currentBest,
        Unit:                    unit,
        ConsecutiveSessionCount: sessions,
    }

    // Validate
    if err := progression.Validate(); err != nil {
        return err
    }

    // Check significance
    if !progression.IsSignificantProgression() {
        return nil  // Not significant, skip
    }

    // Create event
    payload := events.StrengthProgressionPayload{
        ExerciseID:          exerciseID,
        ExerciseName:        exerciseName,
        PreviousBestWeight:  previousBest,
        NewBestWeight:       currentBest,
        Unit:                unit,
        ConsecutiveSessions: sessions,
        DetectedAt:          time.Now(),
    }

    event := events.NewEventEnvelope(
        "StrengthProgressionDetected",
        exerciseID,
        uuid.New().String(),
        payload,
    )

    // Publish event
    return s.eventPublisher.PublishEvent(ctx, event)
}
```

**Key responsibilities:**
- Orchestrate domain logic
- Call ports (interfaces)
- Handle errors
- Publish events

### Step 4: Define Ports (Interfaces)

Ports define how your service communicates with the outside world.

**File:** `backend/services/analytics-service/internal/ports/ports.go`

```go
package ports

import (
    "context"
    "github.com/marcus/ironlog/backend/shared/events"
)

// EventPublisher sends events to Kafka
type EventPublisher interface {
    PublishEvent(ctx context.Context, event *events.EventEnvelope) error
}

// Repository stores/retrieves data
type Repository interface {
    GetExerciseHistory(ctx context.Context, exerciseID string) ([]ExerciseRecord, error)
    SaveProgression(ctx context.Context, progression interface{}) error
}

type ExerciseRecord struct {
    ExerciseID string
    Weight     float64
    Reps       int
    Date       time.Time
}
```

**Why ports?**
- Decouples from implementation details
- Makes testing easy (mock the ports)
- Enables swapping implementations

### Step 5: Implement Adapters (Infrastructure)

Adapters are where technology-specific code lives.

**File:** `backend/services/analytics-service/internal/adapters/outbound/kafka/event_publisher.go`

```go
package kafka

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
    "github.com/marcus/ironlog/backend/shared/events"
)

type EventPublisher struct {
    writer *kafka.Writer
    topic  string
}

func NewEventPublisher(brokers []string, topic string) *EventPublisher {
    writer := &kafka.Writer{
        Addr:     kafka.TCP(brokers...),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
    return &EventPublisher{writer: writer, topic: topic}
}

func (ep *EventPublisher) PublishEvent(
    ctx context.Context,
    event *events.EventEnvelope,
) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }

    message := kafka.Message{
        Key:   []byte(event.AggregateID),
        Value: payload,
    }

    return ep.writer.WriteMessages(ctx, message)
}
```

**Adapter types:**
- **Inbound**: HTTP handlers, gRPC servers
- **Outbound**: Database, Kafka, Redis, external APIs

### Step 6: Connect Everything in main.go

**File:** `backend/services/analytics-service/cmd/main.go`

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/marcus/ironlog/backend/services/analytics-service/internal/adapters/inbound/http"
    kafkaAdapter "github.com/marcus/ironlog/backend/services/analytics-service/internal/adapters/outbound/kafka"
    "github.com/marcus/ironlog/backend/services/analytics-service/internal/adapters/outbound/postgres"
    "github.com/marcus/ironlog/backend/services/analytics-service/internal/application"
    "github.com/marcus/ironlog/backend/shared/infra"
)

func main() {
    // Get config from environment
    kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
    if kafkaBrokers[0] == "" {
        kafkaBrokers = []string{"localhost:9092"}
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    // Initialize adapters
    eventPublisher := kafkaAdapter.NewEventPublisher(kafkaBrokers, "analytics-events")
    repository := postgres.NewRepository()  // Initialize DB

    // Initialize application service
    strengthService := application.NewStrengthProgressionService(
        eventPublisher,
        repository,
    )

    // Initialize HTTP handler
    handler := http.NewStrengthProgressionHandler(strengthService)

    // Setup routes
    http.HandleFunc("/strength-progression", handler.DetectProgression)
    http.HandleFunc("/health", handler.HealthCheck)

    log.Printf("Analytics service starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

**Key points:**
- Initialize adapters first
- Pass adapters to services
- Services don't know about adapters (dependency inversion)

### Step 7: Create HTTP Handler

**File:** `backend/services/analytics-service/internal/adapters/inbound/http/handler.go`

```go
package http

import (
    "encoding/json"
    "net/http"
    "github.com/marcus/ironlog/backend/services/analytics-service/internal/application"
)

type StrengthProgressionHandler struct {
    service *application.StrengthProgressionService
}

func NewStrengthProgressionHandler(
    service *application.StrengthProgressionService,
) *StrengthProgressionHandler {
    return &StrengthProgressionHandler{service: service}
}

type DetectProgressionRequest struct {
    ExerciseID          string  `json:"exercise_id"`
    ExerciseName        string  `json:"exercise_name"`
    PreviousBestWeight  float64 `json:"previous_best_weight"`
    NewBestWeight       float64 `json:"new_best_weight"`
    Unit                string  `json:"unit"`
    ConsecutiveSessions int     `json:"consecutive_sessions"`
}

func (h *StrengthProgressionHandler) DetectProgression(
    w http.ResponseWriter,
    r *http.Request,
) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req DetectProgressionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Call service
    err := h.service.DetectStrengthProgression(
        r.Context(),
        req.ExerciseID,
        req.ExerciseName,
        req.PreviousBestWeight,
        req.NewBestWeight,
        req.Unit,
        req.ConsecutiveSessions,
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *StrengthProgressionHandler) HealthCheck(
    w http.ResponseWriter,
    r *http.Request,
) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
```

### Step 8: Write Tests

**File:** `backend/services/analytics-service/internal/domain/strength_progression_test.go`

```go
package domain

import (
    "testing"
)

func TestIsSignificantProgression(t *testing.T) {
    tests := []struct {
        name     string
        sp       *StrengthProgression
        expected bool
    }{
        {
            name: "Significant: 10% improvement with 3 sessions",
            sp: &StrengthProgression{
                ExerciseID:              "ex_bench",
                PreviousBestWeight:      20.0,
                NewBestWeight:           22.0,  // 10% improvement
                ConsecutiveSessionCount: 3,
            },
            expected: true,
        },
        {
            name: "Not significant: Only 3% improvement",
            sp: &StrengthProgression{
                ExerciseID:              "ex_bench",
                PreviousBestWeight:      20.0,
                NewBestWeight:           20.6,  // 3% improvement
                ConsecutiveSessionCount: 3,
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.sp.IsSignificantProgression()
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

**Test strategy:**
- Test domain logic exhaustively
- Mock ports in application tests
- Integration tests for adapters

### Step 9: Update Frontend (if needed)

**File:** `frontend/web-app/src/services/analytics.ts`

```typescript
interface StrengthProgression {
    exerciseId: string;
    exerciseName: string;
    previousBestWeight: number;
    newBestWeight: number;
    unit: string;
    consecutiveSessions: number;
}

export async function detectStrengthProgression(
    data: StrengthProgression
): Promise<{ status: string }> {
    const response = await fetch(
        'http://localhost:8082/strength-progression',
        {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        }
    );

    if (!response.ok) {
        throw new Error('Failed to detect progression');
    }

    return response.json();
}
```

### Step 10: Test End-to-End

```bash
# 1. Rebuild services
cd backend/services/analytics-service
go build -o analytics-service ./cmd

# 2. Restart container
docker-compose restart analytics

# 3. Test the endpoint
curl -X POST http://localhost:8082/strength-progression \
  -H "Content-Type: application/json" \
  -d '{
    "exercise_id": "ex_bench",
    "exercise_name": "Bench Press",
    "previous_best_weight": 20.0,
    "new_best_weight": 22.0,
    "unit": "kg",
    "consecutive_sessions": 3
  }'

# 4. Check if event was published
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic analytics-events \
  --from-beginning
```

## Debugging Tips

### View Service Logs

```bash
docker-compose logs -f analytics-service
```

### Connect to Database

```bash
docker-compose exec postgres psql -U postgres -d ironlog
SELECT * FROM event_store ORDER BY created_at DESC LIMIT 10;
```

### Monitor Kafka

```bash
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic analytics-events \
  --from-beginning
```

### Check Redis

```bash
docker-compose exec redis redis-cli
KEYS *
GET "event:evt_001"
```

## Adding a New Service

If you need a completely new microservice:

1. **Create directory:** `backend/services/new-service/`
2. **Copy template** from existing service (e.g., `analytics-service`)
3. **Update domain logic** for your use case
4. **Update Dockerfile** in `new-service/`
5. **Add to docker-compose.yml:**
   ```yaml
   new-service:
     build:
       context: ../backend
       dockerfile: ./services/new-service/Dockerfile
     environment:
       KAFKA_BROKERS: kafka:29092
       PORT: 8086
     ports:
       - "8086:8086"
     depends_on:
       - kafka
   ```
6. **Update README** and documentation

## Best Practices

✅ **DO:**
- Start with events
- Keep domain logic pure
- Use interfaces (ports)
- Write tests
- Document decisions

❌ **DON'T:**
- Put business logic in HTTP handlers
- Make domain code depend on frameworks
- Skip tests
- Ignore error handling
- Hardcode configuration

## Next Steps

- Review [SERVICE_COMMUNICATION.md](SERVICE_COMMUNICATION.md) for event flows
- Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues
- Follow [E2E_TEST_GUIDE.md](E2E_TEST_GUIDE.md) to verify your feature
