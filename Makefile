.PHONY: help docker-up docker-down build-services build-frontend run-parser run-workout run-analytics test-parser lint clean

help:
	@echo "IRONLOG Project Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make docker-up           Start all services with Docker Compose"
	@echo "  make docker-down         Stop all services"
	@echo "  make build-services      Build all Go services"
	@echo "  make build-frontend      Build frontend application"
	@echo "  make run-parser          Run parser service locally"
	@echo "  make run-workout         Run workout service locally"
	@echo "  make run-analytics       Run analytics service locally"
	@echo "  make test-parser         Test parser DSL"
	@echo "  make lint                Lint Go code"
	@echo "  make clean               Clean build artifacts"

docker-up:
	cd infra && docker-compose up -d
	@echo "Services starting up..."
	@echo "Frontend: http://localhost:3000"
	@echo "API Gateway: http://localhost:8080"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3001"
	@echo "Jaeger: http://localhost:16686"

docker-down:
	cd infra && docker-compose down

build-services:
	@echo "Building parser-service..."
	cd backend/services/parser-service && go build -o parser-service ./cmd
	@echo "Building workout-service..."
	cd backend/services/workout-service && go build -o workout-service ./cmd
	@echo "Building analytics-service..."
	cd backend/services/analytics-service && go build -o analytics-service ./cmd
	@echo "Building projection-service..."
	cd backend/services/projection-service && go build -o projection-service ./cmd
	@echo "Building recommendation-service..."
	cd backend/services/recommendation-service && go build -o recommendation-service ./cmd
	@echo "Building notification-service..."
	cd backend/services/notification-service && go build -o notification-service ./cmd
	@echo "All services built!"

build-frontend:
	@echo "Building frontend..."
	cd frontend/web-app && npm install && npm run build

run-parser:
	cd backend/services/parser-service && PORT=8081 KAFKA_BROKERS=localhost:9092 go run ./cmd/main.go

run-workout:
	cd backend/services/workout-service && \
	DB_HOST=localhost DB_USER=postgres DB_PASSWORD=password \
	DB_NAME=ironlog KAFKA_BROKERS=localhost:9092 PORT=8080 \
	go run ./cmd/main.go

run-analytics:
	cd backend/services/analytics-service && \
	KAFKA_BROKERS=localhost:9092 PORT=8082 go run ./cmd/main.go

test-parser:
	@echo "Testing parser service..."
	curl -X POST http://localhost:8081/parse \
	  -H "Content-Type: application/json" \
	  -d '{"raw_text": "SUPINO RETO\n\nWarm up: 1x 1-20 10kg (10)\nWork: 3x 6-8 20kg (8 7 6)"}'

test-workout:
	@echo "Starting workout..."
	curl -X POST http://localhost:8080/workouts/start \
	  -H "Content-Type: application/json" \
	  -d '{"workout_id": "wk_001", "user_id": "user_123", "notes": "Test session"}'

lint:
	@echo "Linting Go code..."
	golangci-lint run ./backend/...

clean:
	@echo "Cleaning build artifacts..."
	find backend -name "*.o" -delete
	find backend -name "*.so" -delete
	rm -rf frontend/web-app/dist
	rm -rf frontend/web-app/build
	@echo "Cleanup complete!"
