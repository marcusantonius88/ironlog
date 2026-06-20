# IRONLOG Project Structure

```
ironlog/
│
├── README.md                          # Main documentation
├── QUICKSTART.md                      # Quick start guide
├── Makefile                          # Build automation
├── .env.example                      # Environment template
├── .gitignore                        # Git ignore rules
│
├── backend/
│   ├── services/
│   │   ├── workout-service/          # Core event sourcing service
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   └── workout.go
│   │   │   │   ├── application/
│   │   │   │   │   └── workout_service.go
│   │   │   │   ├── ports/
│   │   │   │   │   └── repositories.go
│   │   │   │   └── adapters/
│   │   │   │       ├── inbound/http/
│   │   │   │       │   └── handler.go
│   │   │   │       └── outbound/postgres/
│   │   │   │           └── event_store.go
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   ├── parser-service/           # DSL parsing service
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   ├── lexer.go      # Tokenization
│   │   │   │   │   └── parser.go     # AST generation
│   │   │   │   ├── application/
│   │   │   │   │   └── parsing_service.go
│   │   │   │   └── adapters/
│   │   │   │       ├── inbound/http/
│   │   │   │       │   └── handler.go
│   │   │   │       └── outbound/kafka/
│   │   │   │           └── event_publisher.go
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   ├── analytics-service/        # Performance analytics
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/domain/
│   │   │   │   └── analytics.go
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   ├── projection-service/       # CQRS read models
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/domain/
│   │   │   │   └── projections.go
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   ├── recommendation-service/   # Recommendations engine
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/domain/
│   │   │   │   └── recommendations.go
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   └── notification-service/     # Notifications
│   │       ├── cmd/main.go
│   │       ├── internal/domain/
│   │       │   └── notifications.go
│   │       ├── Dockerfile
│   │       └── go.mod
│   │
│   └── shared/                       # Shared code & contracts
│       ├── events/
│       │   ├── event_envelope.go     # Base event structure
│       │   ├── workout_events.go
│       │   ├── analytics_events.go
│       │   └── parser_events.go
│       ├── infra/
│       │   ├── kafka.go              # Kafka producer/consumer
│       │   ├── database.go           # PostgreSQL connection
│       │   ├── redis.go              # Redis client
│       │   └── observability.go      # Metrics & tracing
│       ├── database/
│       │   └── migrations.go         # SQL migrations
│       ├── go.mod
│       └── go.sum
│
├── frontend/
│   └── web-app/
│       ├── src/
│       │   ├── components/
│       │   │   ├── Navigation.tsx
│       │   │   └── DSLParser.tsx
│       │   ├── pages/
│       │   │   ├── Dashboard.tsx
│       │   │   ├── WorkoutEntry.tsx
│       │   │   └── Analytics.tsx
│       │   ├── services/
│       │   ├── hooks/
│       │   ├── types/
│       │   ├── App.tsx
│       │   ├── main.tsx
│       │   └── index.css
│       ├── index.html
│       ├── package.json
│       ├── vite.config.ts
│       ├── tailwind.config.ts
│       ├── tsconfig.json
│       ├── Dockerfile
│       └── .gitignore
│
├── infra/
│   ├── docker-compose.yml            # Complete infrastructure
│   ├── postgres/
│   │   └── init.sql                  # Database schema
│   ├── prometheus/
│   │   └── prometheus.yml            # Metrics configuration
│   ├── grafana/
│   │   ├── provisioning/
│   │   │   ├── datasources/
│   │   │   │   └── prometheus.yml
│   │   │   └── dashboards/
│   │   │       └── config.yml
│   ├── kafka/                        # Kafka configuration
│   ├── debezium/                     # CDC configuration
│   ├── redis/                        # Redis configuration
│   └── otel/                         # OpenTelemetry configuration
│
└── docs/
    ├── ARCHITECTURE.md
    ├── DSL_SPEC.md
    ├── EVENT_CONTRACTS.md
    ├── DEPLOYMENT.md
    └── TROUBLESHOOTING.md
```

## Key Design Decisions

### Event Sourcing
- Complete audit trail of all changes
- Ability to replay history
- Temporal analysis capabilities

### CQRS Separation
- Optimized write model (event store)
- Optimized read models (projections)
- Scalable independent scaling

### Microservices
- Each service handles single responsibility
- Loose coupling via events
- Easy to develop and test independently

### Clean Architecture
- Domain layer independent of frameworks
- Ports for interfaces/contracts
- Adapters for infrastructure concerns

### DSL Parser
- Custom lexer for tokenization
- Deterministic AST parser
- No AI-based parsing (explicit rules)

## Technology Selection Rationale

| Component | Choice | Reason |
|-----------|--------|--------|
| Backend | Go | Performance, simplicity, concurrency |
| Event Store | PostgreSQL | ACID, JSONB, full-text search |
| Streaming | Kafka | High throughput, replay capability |
| CDC | Debezium | Minimal code, proven reliability |
| Caching | Redis | Fast, simple, idempotency support |
| Metrics | Prometheus | Standard, Grafana integration |
| Frontend | React | Modern, component-based, rich ecosystem |
| Build | Vite | Fast, modern, ESM-first |

## Service Responsibilities

| Service | Responsibility | Events Consumed | Events Published |
|---------|-----------------|-----------------|------------------|
| **parser** | DSL parsing, AST generation | ParsingRequested | DSLParsedSuccessfully, DSLParsingFailed |
| **workout** | Event store, aggregates, commands | - | WorkoutStarted, SetPerformed, WorkoutFinished |
| **analytics** | Metrics, trends, PR detection | SetPerformed, WorkoutFinished | PerformanceImproved, PersonalRecordAchieved |
| **projection** | Read models, CQRS | All domain events | ProjectionUpdated |
| **recommendation** | Suggestions, load progression | PerformanceImproved | LoadIncreaseSuggested |
| **notification** | Alert dispatch | All analytics events | NotificationSent |

