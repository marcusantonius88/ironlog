# Docker Deployment Guide - IRONLOG

## Overview

IRONLOG is deployed using **Docker Compose**, which orchestrates all services and infrastructure components in a single configuration.

## Quick Start

```bash
cd infra
docker-compose up -d
```

This starts all 13 components:
- PostgreSQL (event store)
- Kafka + Zookeeper (event streaming)
- Debezium (CDC)
- Redis (caching & idempotency)
- Prometheus (metrics)
- Grafana (dashboards)
- Jaeger (distributed tracing)
- 6 Backend services (parser, workout, analytics, projection, recommendation, notification)
- Frontend React app

## System Requirements

- **Docker Engine**: 20.10+
- **Docker Compose**: 2.0+
- **Disk Space**: 10GB minimum (for volumes)
- **RAM**: 8GB recommended (4GB minimum)
- **CPU**: 4 cores recommended

## Configuration

### Environment Variables

Each service reads its configuration from environment variables defined in `docker-compose.yml`:

**Database (PostgreSQL):**
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=ironlog
```

**Kafka:**
```env
KAFKA_BROKER_ID=1
KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
```

**Services:**
```env
# Parser Service
KAFKA_BROKERS=kafka:29092
PORT=8081

# Workout Service
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=ironlog
KAFKA_BROKERS=kafka:29092
PORT=8080
```

### Custom Configuration

To modify settings, edit `infra/docker-compose.yml`:

```yaml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: your_password  # Change here
      
  parser-service:
    environment:
      KAFKA_BROKERS: your-kafka:9092    # Change here
```

Then restart:
```bash
docker-compose restart
```

## Service Ports

| Service | Port | Purpose |
|---------|------|---------|
| **Workout Service** | 8080 | Main API |
| **Parser Service** | 8081 | DSL Parsing API |
| **Analytics Service** | 8082 | Analytics API |
| **Projection Service** | 8083 | CQRS Read Models |
| **Recommendation Service** | 8084 | Recommendations API |
| **Notification Service** | 8085 | Notifications API |
| **Frontend** | 3000 | React Web App |
| **PostgreSQL** | 5432 | Event Store |
| **Kafka** | 9092 | Event Streaming |
| **Redis** | 6379 | Caching |
| **Prometheus** | 9090 | Metrics |
| **Grafana** | 3001 | Dashboards |
| **Jaeger** | 16686 | Distributed Tracing |
| **Debezium** | 8083 | CDC Server |

## Operations

### Check Status

```bash
# List all running containers
docker-compose ps

# View logs
docker-compose logs -f                    # All services
docker-compose logs -f workout-service    # Single service
docker-compose logs -f postgres          # Specific infrastructure
```

### Restart Services

```bash
# Restart specific service
docker-compose restart parser-service

# Restart all services
docker-compose restart

# Restart and rebuild (if code changed)
docker-compose up -d --build
```

### Stop and Cleanup

```bash
# Stop all containers (keep volumes)
docker-compose stop

# Stop and remove containers (keep volumes)
docker-compose down

# Stop, remove containers AND volumes (hard reset)
docker-compose down -v

# Remove specific volume
docker volume rm ironlog_postgres_data
```

### View Database

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d ironlog

# Common queries
\dt                                    # List tables
\d event_store                         # Describe table
SELECT COUNT(*) FROM event_store;      # Row count
SELECT * FROM event_store LIMIT 10;    # View data
\q                                     # Exit
```

### View Kafka Topics

```bash
# List topics
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Monitor messages in real-time
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic workout-events \
  --from-beginning

# Create topic
docker-compose exec kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --partitions 1 \
  --replication-factor 1
```

### View Redis Keys

```bash
# Connect to Redis
docker-compose exec redis redis-cli

# Common commands
PING              # Test connection
SET key value     # Store value
GET key           # Retrieve value
KEYS *            # List all keys
DEL key           # Delete key
FLUSHDB           # Clear all keys
EXIT              # Exit
```

## Health Checks

All services expose `/health` endpoints:

```bash
# Backend services
curl http://localhost:8080/health   # Workout Service
curl http://localhost:8081/health   # Parser Service
curl http://localhost:8082/health   # Analytics Service
curl http://localhost:8083/health   # Projection Service
curl http://localhost:8084/health   # Recommendation Service
curl http://localhost:8085/health   # Notification Service

# Infrastructure
curl http://localhost:9090/-/healthy  # Prometheus
curl http://localhost:6379 PING       # Redis (via redis-cli)
```

Expected response:
```json
{"status":"healthy"}
```

## Monitoring & Observability

### Prometheus Metrics

Access: **http://localhost:9090**

Query examples:
```promql
# All metrics
up

# Service-specific
up{job="parser-service"}
up{job="workout-service"}

# Request count
parsing_requests_total
events_processed_total
```

### Grafana Dashboards

Access: **http://localhost:3001**
- **Username**: admin
- **Password**: password

Dashboards are auto-provisioned:
- Service Health
- Event Processing
- Database Performance
- Kafka Topics

### Jaeger Distributed Tracing

Access: **http://localhost:16686**

1. Select service from dropdown
2. Click "Find Traces"
3. View distributed call chains

## Backup & Recovery

### Backup Database

```bash
# Create backup
docker-compose exec postgres pg_dump -U postgres ironlog > backup.sql

# Restore from backup
docker-compose exec -T postgres psql -U postgres ironlog < backup.sql
```

### Backup Volumes

```bash
# List volumes
docker volume ls | grep ironlog

# Backup volume data
docker run --rm -v ironlog_postgres_data:/data \
  -v $(pwd):/backup alpine tar czf /backup/postgres.tar.gz -C /data .

# Restore volume data
docker run --rm -v ironlog_postgres_data:/data \
  -v $(pwd):/backup alpine tar xzf /backup/postgres.tar.gz -C /data
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in docker-compose.yml
ports:
  - "8090:8080"  # Change host port
```

### Services Not Starting

```bash
# Check logs
docker-compose logs workout-service

# Rebuild images
docker-compose down
docker-compose up -d --build

# Check Docker resources
docker stats
```

### Database Connection Errors

```bash
# Verify PostgreSQL is running
docker-compose logs postgres

# Test connection
docker-compose exec postgres psql -U postgres -c "SELECT version();"

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### Kafka Connection Issues

```bash
# Check Kafka status
docker-compose logs kafka

# Verify broker is ready
docker-compose exec kafka kafka-broker-api-versions --bootstrap-server localhost:29092

# Check topic creation
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:29092
```

## Performance Tuning

### Increase Resource Limits

Edit `docker-compose.yml`:

```yaml
services:
  postgres:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G

  kafka:
    environment:
      KAFKA_HEAP_OPTS: "-Xmx2G -Xms2G"
```

### Database Optimization

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d ironlog

# Analyze table
ANALYZE event_store;

# Vacuum table
VACUUM ANALYZE event_store;

# Create index
CREATE INDEX idx_aggregate_id ON event_store(aggregate_id);

# View indexes
\di

\q
```

## Next Steps

- Read [QUICKSTART.md](QUICKSTART.md) for first-time setup
- Check [E2E_TEST_GUIDE.md](E2E_TEST_GUIDE.md) to test all components
- Review [ARCHITECTURE.md](ARCHITECTURE.md) to understand system design
