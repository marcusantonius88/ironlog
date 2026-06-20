# Quick Start Guide for IRONLOG

## Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- Git

## Option 1: Full Docker Setup (Recommended)

```bash
# Clone and navigate
cd ironlog

# Run setup script
chmod +x setup.sh
./setup.sh

# Or manually:
cd infra
docker-compose up -d
cd ..

# Start frontend development server
cd frontend/web-app
npm install
npm run dev
```

## Option 2: Local Development

### 1. Start Infrastructure

```bash
cd infra
docker-compose up -d postgres kafka zookeeper redis prometheus grafana jaeger
```

### 2. Run Services

Terminal 1 - Parser Service:
```bash
cd backend/services/parser-service
PORT=8081 KAFKA_BROKERS=localhost:9092 go run ./cmd/main.go
```

Terminal 2 - Workout Service:
```bash
cd backend/services/workout-service
DB_HOST=localhost DB_USER=postgres DB_PASSWORD=password \
DB_NAME=ironlog KAFKA_BROKERS=localhost:9092 PORT=8080 \
go run ./cmd/main.go
```

Terminal 3 - Frontend:
```bash
cd frontend/web-app
npm install
npm run dev
```

### 3. Access Services

- Frontend: http://localhost:3000
- Workout API: http://localhost:8080
- Parser API: http://localhost:8081
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001
- Jaeger: http://localhost:16686

## Testing DSL Parser

```bash
curl -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d '{
    "raw_text": "SUPINO RETO BARRA\n\nWarm up: 1x 1-20 10kg (10)\nFeeder: 1x 08-10 15kg (10)\nWork: 2x 08-10 20kg (6 8)\nTop set: 1x 6-8 25kg (3)"
  }'
```

Expected response:
```json
{
  "success": true,
  "exercise": "SUPINO RETO BARRA",
  "set_groups": [
    {
      "set_type": "WARM_UP",
      "planned_sets": 1,
      "target_rep_min": 1,
      "target_rep_max": 20,
      "weight": 10.0,
      "unit": "kg",
      "executed_reps": [10]
    },
    ...
  ]
}
```

## Recording a Workout

### Start Workout
```bash
curl -X POST http://localhost:8080/workouts/start \
  -H "Content-Type: application/json" \
  -d '{
    "workout_id": "wk_20240512_001",
    "user_id": "user_marcus",
    "notes": "Heavy day"
  }'
```

### Record Set
```bash
curl -X POST http://localhost:8080/workouts/sets/perform \
  -H "Content-Type: application/json" \
  -d '{
    "workout_id": "wk_20240512_001",
    "exercise_id": "ex_supino_001",
    "set_number": 1,
    "set_type": "WORK",
    "planned_sets": 3,
    "planned_weight": 20.0,
    "executed_reps": [8, 7, 6],
    "target_rep_min": 6,
    "target_rep_max": 8,
    "notes": "Good form"
  }'
```

### Finish Workout
```bash
curl -X POST http://localhost:8080/workouts/finish \
  -H "Content-Type: application/json" \
  -d '{
    "workout_id": "wk_20240512_001",
    "notes": "Great session!"
  }'
```

## Useful Commands

```bash
# Build all services
make build-services

# Build frontend
make build-frontend

# Start Docker infrastructure
make docker-up

# Stop Docker infrastructure
make docker-down

# Run individual services locally
make run-parser
make run-workout
make run-analytics

# Test parser
make test-parser

# Clean build artifacts
make clean
```

## Monitoring & Observability

### Prometheus Metrics
- Visit: http://localhost:9090
- Query: `events_processed_total`, `parsing_requests_total`

### Grafana Dashboards
- Visit: http://localhost:3001
- Login: admin/password
- Pre-configured Prometheus datasource

### Jaeger Tracing
- Visit: http://localhost:16686
- Search traces by service name
- See distributed traces across services

### Health Checks
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## Debugging

### View Kafka Topics
```bash
docker exec ironlog-kafka kafka-topics --list --bootstrap-server localhost:9092
```

### View PostgreSQL Events
```bash
docker exec -it ironlog-postgres psql -U postgres -d ironlog -c "SELECT * FROM event_store ORDER BY created_at DESC LIMIT 10;"
```

### View Redis Keys
```bash
docker exec -it ironlog-redis redis-cli KEYS "*"
```

### Service Logs
```bash
# Docker services
docker-compose logs -f parser-service
docker-compose logs -f workout-service

# Local services (check terminal output)
```

## Troubleshooting

### Port Already in Use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or change port in .env
```

### Database Connection Error
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check credentials in .env
cat .env | grep DB_
```

### Kafka Connection Error
```bash
# Check Kafka is running
docker ps | grep kafka

# Test connectivity
docker exec ironlog-kafka kafka-broker-api-versions --bootstrap-server localhost:9092
```

## Next Steps

1. Review architecture in [README.md](../README.md)
2. Explore DSL syntax in [DSL_SPEC.md](../docs/DSL_SPEC.md)
3. Study event model in `backend/shared/events/`
4. Check service implementations in `backend/services/`

Happy training! 🏋️
