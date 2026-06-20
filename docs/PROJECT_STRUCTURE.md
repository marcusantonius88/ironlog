# Project Structure Guide - IRONLOG

## Directory Overview

```
ironlog/
├── backend/
│   ├── services/
│   │   ├── parser-service/        # DSL parsing engine
│   │   ├── workout-service/       # Core event sourcing service
│   │   ├── analytics-service/     # Performance metrics & insights
│   │   ├── projection-service/    # CQRS read model builder
│   │   ├── recommendation-service/# Load progression suggestions
│   │   └── notification-service/  # Alert dispatcher
│   └── shared/
│       ├── events/                # Shared event definitions
│       ├── infra/                 # Infrastructure utilities
│       ├── database/              # Database utilities & migrations
│       └── go.mod                 # Shared Go module
│
├── frontend/
│   └── web-app/                   # React TypeScript application
│       ├── src/
│       │   ├── components/        # Reusable UI components
│       │   ├── pages/             # Page-level components
│       │   ├── hooks/             # Custom React hooks
│       │   ├── services/          # API clients & utilities
│       │   ├── types/             # TypeScript interfaces
│       │   ├── utils/             # Helper functions
│       │   ├── App.tsx            # Root component
│       │   └── main.tsx           # Entry point
│       └── vite.config.ts         # Vite configuration
│
├── infra/
│   ├── docker-compose.yml         # Infrastructure orchestration
│   ├── postgres/
│   │   └── init.sql               # Database schema & migrations
│   ├── prometheus/
│   │   └── prometheus.yml         # Metrics scraping config
│   ├── grafana/
│   │   └── provisioning/          # Pre-configured dashboards
│   ├── kafka/                     # Kafka configuration
│   ├── debezium/                  # CDC connectors
│   ├── redis/                     # Redis persistence config
│   └── otel/                      # OpenTelemetry tracing
│
├── docs/
│   ├── QUICKSTART.md              # Getting started guide
│   ├── ARCHITECTURE.md            # System design & patterns
│   ├── DSL_SPEC.md                # DSL language specification
│   ├── SERVICE_COMMUNICATION.md   # Event flow & contracts
│   ├── PROJECT_STRUCTURE.md       # This file (directory guide)
│   ├── DEVELOPMENT_GUIDE.md       # How to develop features
│   ├── E2E_TEST_GUIDE.md          # Manual testing procedures
│   ├── DEPLOYMENT_DOCKER.md       # Docker deployment guide
│   ├── TROUBLESHOOTING.md         # Common issues & solutions
│   └── banner-ironlog.png         # Project logo
│
├── setup.sh                       # Initial setup script
├── init.sh                        # Initialization script
├── Makefile                       # Build automation
├── README.md                      # Project overview
└── .gitignore                     # Git ignore rules
```

## Backend Services Structure

Each microservice follows **Hexagonal Architecture** (Ports & Adapters):

```
parser-service/
├── cmd/
│   └── main.go                    # Service entry point & HTTP server setup
├── internal/
│   ├── domain/                    # Business logic (pure Go, framework-agnostic)
│   │   ├── lexer.go               # Token generation
│   │   ├── parser.go              # AST parsing
│   │   └── models.go              # Domain entities
│   ├── application/               # Use cases & orchestration
│   │   ├── parsing_service.go     # Main parsing use case
│   │   └── dto.go                 # Data transfer objects
│   ├── ports/                     # Interfaces (contracts)
│   │   ├── repositories.go        # Storage interface
│   │   └── event_publisher.go     # Event publishing interface
│   └── adapters/                  # Infrastructure implementations
│       ├── inbound/
│       │   └── http/
│       │       ├── handler.go     # HTTP request handlers
│       │       └── middleware.go  # HTTP middleware
│       └── outbound/
│           ├── kafka/
│           │   └── event_publisher.go   # Kafka producer
│           └── postgres/
│               └── repository.go         # PostgreSQL implementation
├── Dockerfile                     # Service container image
├── go.mod                         # Go module definition
└── go.sum                         # Go dependencies lock file
```

### Service Template Breakdown

**Domain Layer** (`internal/domain/`)
- Contains pure business logic
- No external dependencies
- Fully testable in isolation

**Application Layer** (`internal/application/`)
- Orchestrates domain logic
- Handles use cases
- Coordinates ports (interfaces)

**Ports** (`internal/ports/`)
- Interface definitions
- Service contracts
- Technology-agnostic

**Adapters** (`internal/adapters/`)
- HTTP handlers (inbound)
- Database implementations (outbound)
- Kafka/Redis clients (outbound)
- Technology-specific implementations

## Shared Code Structure

```
backend/shared/
├── events/
│   ├── event_envelope.go           # Base event type
│   ├── workout_events.go           # Workout event payloads
│   ├── analytics_events.go         # Analytics event payloads
│   ├── parser_events.go            # Parser event payloads
│   └── recommendation_events.go    # Recommendation event payloads
│
├── infra/
│   ├── kafka.go                    # Kafka producer/consumer
│   ├── database.go                 # PostgreSQL connection pool
│   ├── redis.go                    # Redis client wrapper
│   └── observability.go            # Prometheus metrics & tracing
│
├── database/
│   └── migrations.go               # Database schema setup
│
└── go.mod                          # Shared module
```

### Event Definitions

All events published to Kafka are defined here, ensuring consistency across services:

```go
// workout_events.go
type WorkoutStartedPayload struct {
    WorkoutID string
    UserID    string
    StartedAt time.Time
    Notes     string
}

type SetPerformedPayload struct {
    WorkoutID    string
    ExerciseID   string
    SetNumber    int
    ExecutedReps []int
    // ... more fields
}
```

## Frontend Structure

```
frontend/web-app/
├── src/
│   ├── components/
│   │   ├── DSLParser.tsx           # DSL input & parsing UI
│   │   ├── Navigation.tsx          # Main navigation
│   │   ├── Dashboard.tsx           # Analytics dashboard
│   │   └── WorkoutRecorder.tsx     # Workout logging UI
│   │
│   ├── pages/
│   │   ├── WorkoutEntry.tsx        # Workout recording page
│   │   ├── Dashboard.tsx           # Analytics view page
│   │   └── Settings.tsx            # User preferences page
│   │
│   ├── hooks/
│   │   ├── useWorkout.ts           # Workout state management
│   │   ├── useParser.ts            # DSL parsing hook
│   │   └── useFetch.ts             # HTTP request utility
│   │
│   ├── services/
│   │   ├── api.ts                  # Backend API client
│   │   ├── parser.ts               # Parser service wrapper
│   │   └── workout.ts              # Workout service wrapper
│   │
│   ├── types/
│   │   ├── workout.ts              # Workout interfaces
│   │   ├── api.ts                  # API response types
│   │   └── events.ts               # Event interfaces
│   │
│   ├── utils/
│   │   ├── validators.ts           # Input validation
│   │   ├── formatters.ts           # Data formatting
│   │   └── constants.ts            # App constants
│   │
│   ├── App.tsx                     # Root component
│   ├── main.tsx                    # Vite entry point
│   └── index.css                   # Global styles
│
├── package.json                    # Dependencies & scripts
├── tsconfig.json                   # TypeScript config
├── vite.config.ts                  # Vite bundler config
├── tailwind.config.ts              # Tailwind CSS config
├── postcss.config.cjs              # PostCSS config
└── index.html                      # HTML template
```

## Infrastructure Configuration

### Docker Compose Services

All services are defined in `infra/docker-compose.yml`:

**Storage:**
- `postgres` - Event store & read models
- `redis` - Caching & idempotency

**Messaging:**
- `kafka` - Event streaming
- `zookeeper` - Kafka coordination
- `debezium` - CDC connector

**Backend Services:**
- `parser-service` through `notification-service`

**Observability:**
- `prometheus` - Metrics collection
- `grafana` - Visualization dashboards
- `jaeger` - Distributed tracing

**Frontend:**
- `web-app` - React application

### Database Schema

Located in `infra/postgres/init.sql`:

```sql
-- Event Store (immutable)
CREATE TABLE event_store (
    event_id UUID PRIMARY KEY,
    event_type VARCHAR NOT NULL,
    aggregate_id VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- CQRS Read Models
CREATE TABLE exercise_progression (
    exercise_id VARCHAR PRIMARY KEY,
    exercise_name VARCHAR,
    current_weight DECIMAL,
    current_reps INT,
    last_updated TIMESTAMP
);

CREATE TABLE personal_records (
    exercise_id VARCHAR PRIMARY KEY,
    max_weight DECIMAL,
    max_reps INT,
    achieved_date DATE
);
```

## Naming Conventions

### Go Services

- **Package names**: lowercase, no underscores (`application`, `adapters`)
- **Exported functions**: PascalCase (`NewParsingService`, `ParseDSL`)
- **Unexported functions**: camelCase (`parseTokens`, `buildAST`)
- **Interface names**: end with `-er` (`Reader`, `Writer`, `Publisher`)

### Go Files

- **Domain logic**: `domain.go`, `models.go`, `lexer.go`, `parser.go`
- **Use cases**: `service.go`, `handler.go`
- **Data transfer**: `dto.go`
- **Implementations**: `postgres.go`, `kafka.go`, `http.go`

### React Components

- **File names**: PascalCase for components (`DSLParser.tsx`)
- **Export names**: Match filename (`export function DSLParser()`)
- **Custom hooks**: `use` prefix (`useWorkout`, `useParser`)
- **Interfaces**: `I` prefix in types file (`IWorkout`, `IEvent`)

### Database Tables

- **Table names**: snake_case, plural (`event_store`, `personal_records`)
- **Column names**: snake_case (`user_id`, `created_at`)
- **Primary keys**: `*_id` pattern
- **Indexes**: `idx_*` pattern

### Identifiers

- **Aggregate IDs**: `wk_*` (workout), `ex_*` (exercise), `usr_*` (user)
- **Event IDs**: `evt_*` or generated UUID
- **Correlation IDs**: `corr_*` for request tracing
- **Request IDs**: `req_*` for parsing requests

## Configuration Files

### Environment Configuration

`.env.example` (in root):
```env
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password

# Kafka
KAFKA_BROKERS=localhost:9092

# Services
PARSER_PORT=8081
WORKOUT_PORT=8080
```

### Build Configuration

**Go Modules**: Each service has `go.mod` defining dependencies
**Frontend**: `package.json` with npm scripts
**Docker**: Each service has `Dockerfile` for containerization

## File Dependencies

```
                          shared/events/
                               ↑
    ┌───────────────────────────┼───────────────────────┐
    │                           │                       │
parser-service        workout-service            analytics-service
    │                           │                       │
    ├─ Consumes: (HTTP)         ├─ Publishes: workout   ├─ Consumes: workout
    ├─ Publishes: parser-events │  events              └─ Publishes: analytics
    │                           │                          events
    └─ Depends: shared/infra    └─ Depends: shared/*

                          shared/infra/
                               ↑
                    ┌───────────┼───────────┐
                    │           │           │
            postgres.go   kafka.go    redis.go
                    │           │           │
                    ▼           ▼           ▼
              Event Store    Kafka      Redis Cache
```

## Development Workflow

### Adding a New Feature

1. **Update events** in `backend/shared/events/`
2. **Create domain logic** in service's `domain/` folder
3. **Implement use case** in service's `application/`
4. **Add HTTP handler** in service's `adapters/inbound/http/`
5. **Publish events** in service's `adapters/outbound/kafka/`
6. **Add frontend UI** in `frontend/web-app/src/components/`
7. **Test end-to-end** with [E2E_TEST_GUIDE.md](E2E_TEST_GUIDE.md)

See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) for detailed instructions.

## Next Steps

- Review [ARCHITECTURE.md](ARCHITECTURE.md) for design patterns
- Check [SERVICE_COMMUNICATION.md](SERVICE_COMMUNICATION.md) for event flows
- See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) to add new features
