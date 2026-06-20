# Troubleshooting Guide - IRONLOG

## Container Issues

### Containers Not Starting

**Symptom:** Services fail to start or keep restarting

**Diagnosis:**
```bash
docker-compose logs <service_name>
docker-compose ps  # Check status
```

**Solutions:**

**1. Insufficient Disk Space**
```bash
docker system df  # Check disk usage

# Clean up
docker system prune -a  # Remove unused images/containers
docker volume prune      # Remove unused volumes
```

**2. Port Already in Use**
```bash
# Find what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in docker-compose.yml
ports:
  - "8090:8080"  # Map to different port
```

**3. Memory Issues**
```bash
docker stats  # Check memory usage

# Increase Docker memory limit in preferences
# Or reduce services: docker-compose up postgres kafka parser-service
```

**4. Service Dependencies Not Ready**
```bash
# Check if dependency is up
docker-compose logs postgres

# Increase wait time in depends_on (if custom health check available)
# Or manually wait and restart
sleep 30 && docker-compose restart parser-service
```

### Service Crashes on Startup

**Symptom:** "Container exited with code 1" or "Container suddenly stopped"

**Common Causes:**

**1. Missing Environment Variables**
```bash
# Check what env vars the service needs
docker-compose logs workout-service

# Example error: "DB_HOST not set"
# Solution: Update docker-compose.yml environment section
```

**2. Cannot Connect to Database**
```bash
# Test PostgreSQL connection
docker-compose exec postgres psql -U postgres -c "SELECT 1"

# If fails, reset PostgreSQL
docker-compose stop postgres
docker volume rm ironlog_postgres_data
docker-compose start postgres
sleep 30  # Wait for initialization
```

**3. Cannot Connect to Kafka**
```bash
# Test Kafka connectivity
docker-compose exec kafka kafka-broker-api-versions \
  --bootstrap-server localhost:29092

# If fails, restart Kafka and Zookeeper
docker-compose restart zookeeper kafka
sleep 30
```

---

## Database Issues

### Cannot Connect to PostgreSQL

**Symptom:** "connection refused" or "cannot connect to database"

**Diagnosis:**
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Test connection
docker-compose exec postgres psql -U postgres -c "SELECT version();"
```

**Solutions:**

**1. PostgreSQL Not Started**
```bash
docker-compose start postgres
sleep 30  # Wait for initialization
```

**2. Wrong Credentials**
```bash
# Default credentials (from docker-compose.yml)
# User: postgres
# Password: password
# Database: ironlog

# Connect with correct credentials
docker-compose exec postgres psql -U postgres -d ironlog
```

**3. Port Binding Failed**
```bash
# Check what's using port 5432
lsof -i :5432

# Use different port
# In docker-compose.yml:
ports:
  - "5433:5432"  # Map to 5433 instead
```

### Database Migrations Failed

**Symptom:** Missing tables or schema errors

**Diagnosis:**
```bash
# Check available tables
docker-compose exec postgres psql -U postgres -d ironlog -c "\dt"

# Expected tables:
# event_store, exercise_progression, personal_records, etc.
```

**Solutions:**

**1. Init Script Didn't Run**
```bash
# Manually run migrations
docker-compose exec postgres psql -U postgres -d ironlog < \
  ../infra/postgres/init.sql
```

**2. Database Already Exists (old schema)**
```bash
# Delete and recreate
docker-compose down -v  # Remove volumes
docker-compose up -d postgres
sleep 30  # Wait for init script

# Verify tables
docker-compose exec postgres psql -U postgres -d ironlog -c "\dt"
```

### Query Returning Empty Results

**Symptom:** `SELECT * FROM event_store` returns 0 rows

**Diagnosis:**
```bash
# Check if events are being published
docker-compose logs parser-service | grep "published"

# Check Kafka topics
docker-compose exec kafka kafka-topics --list \
  --bootstrap-server localhost:9092
```

**Solutions:**

**1. Events Not Being Published**
```bash
# Trigger an event (e.g., parse DSL)
curl -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d '{"raw_text": "BENCH\nWork: 1x 8-10 20kg"}'

# Wait 5 seconds for event processing

# Check again
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "SELECT COUNT(*) FROM event_store;"
```

**2. Read Model Not Synchronized**
```bash
# CQRS read models are eventually consistent
# They may lag behind the event store

# Check event store directly
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "SELECT * FROM event_store LIMIT 5;"

# Check projection service logs
docker-compose logs projection-service | tail -50
```

---

## Kafka Issues

### Cannot Publish/Consume Events

**Symptom:** "Failed to publish event" or "Timeout connecting to broker"

**Diagnosis:**
```bash
# Check Kafka is running
docker-compose ps kafka zookeeper

# Check logs
docker-compose logs kafka zookeeper

# Test connectivity
docker-compose exec kafka kafka-broker-api-versions \
  --bootstrap-server localhost:29092
```

**Solutions:**

**1. Zookeeper Not Ready**
```bash
# Kafka depends on Zookeeper; restart Zookeeper first
docker-compose restart zookeeper
sleep 10

# Then restart Kafka
docker-compose restart kafka
sleep 20
```

**2. Broker Not Fully Initialized**
```bash
# Wait longer and check status
sleep 60
docker-compose exec kafka kafka-broker-api-versions \
  --bootstrap-server localhost:29092

# If still fails, hard restart
docker-compose down
docker-compose up -d kafka zookeeper
sleep 60
```

**3. Topic Doesn't Exist**
```bash
# Create topic manually
docker-compose exec kafka kafka-topics --create \
  --bootstrap-server localhost:29092 \
  --topic my-topic \
  --partitions 1 \
  --replication-factor 1

# Verify
docker-compose exec kafka kafka-topics --list \
  --bootstrap-server localhost:29092
```

### Consumer Lag/Messages Not Flowing

**Symptom:** Messages published but not consumed by services

**Diagnosis:**
```bash
# Check consumer groups
docker-compose exec kafka kafka-consumer-groups --list \
  --bootstrap-server localhost:29092

# Check group lag
docker-compose exec kafka kafka-consumer-groups --describe \
  --bootstrap-server localhost:29092 \
  --group parser-service-group

# Monitor messages in real-time
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:29092 \
  --topic workout-events \
  --from-beginning \
  --max-messages 10
```

**Solutions:**

**1. Consumer Service Not Running**
```bash
# Restart the consumer service
docker-compose restart analytics-service

# Check logs
docker-compose logs -f analytics-service
```

**2. Consumer Group Stuck**
```bash
# Reset consumer group offset to earliest
docker-compose exec kafka kafka-consumer-groups --reset-offsets \
  --bootstrap-server localhost:29092 \
  --group analytics-service-group \
  --to-earliest \
  --execute

# Restart consumer
docker-compose restart analytics-service
```

**3. Message Format Invalid**
```bash
# Check message contents
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:29092 \
  --topic workout-events \
  --from-beginning \
  --max-messages 1

# Check service logs for parsing errors
docker-compose logs analytics-service | grep -i "error"
```

---

## Service Communication Issues

### API Endpoint Returns 404

**Symptom:** `curl http://localhost:8080/health` → "Connection refused" or "404 Not Found"

**Diagnosis:**
```bash
# Check if service is running
docker-compose ps workout-service

# Check service logs
docker-compose logs workout-service

# Check if listening on correct port
docker-compose exec workout-service netstat -tlnp | grep 8080
```

**Solutions:**

**1. Service Not Running**
```bash
docker-compose start workout-service
sleep 10
curl http://localhost:8080/health
```

**2. Wrong Port Number**
```bash
# Check docker-compose.yml for port mapping
grep -A 5 "workout-service" infra/docker-compose.yml

# Expected: ports: "8080:8080"
```

**3. Service Crashed**
```bash
# Check logs
docker-compose logs workout-service --tail 100

# Restart with rebuild
docker-compose up -d --build workout-service
```

### Slow Response Times

**Symptom:** Requests taking 10+ seconds to respond

**Diagnosis:**
```bash
# Check service resource usage
docker stats --no-stream workout-service

# Check database query performance
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "SELECT * FROM pg_stat_statements WHERE query LIKE '%workout%' \
      ORDER BY mean_exec_time DESC LIMIT 5;"

# Check service logs
docker-compose logs workout-service
```

**Solutions:**

**1. Database Slow Queries**
```bash
# Analyze tables
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "ANALYZE event_store; ANALYZE exercise_progression;"

# Vacuum (cleanup dead tuples)
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "VACUUM ANALYZE event_store;"

# Create indexes if missing
docker-compose exec postgres psql -U postgres -d ironlog \
  -c "CREATE INDEX IF NOT EXISTS idx_aggregate_id 
      ON event_store(aggregate_id);"
```

**2. Memory Pressure**
```bash
# Check memory
docker stats workout-service

# Restart service
docker-compose restart workout-service

# Or increase Docker memory limit
```

**3. Kafka Bottleneck**
```bash
# Check Kafka broker logs
docker-compose logs kafka

# Check consumer lag
docker-compose exec kafka kafka-consumer-groups --describe \
  --bootstrap-server localhost:29092 \
  --group analytics-service-group
```

---

## Observability Issues

### Prometheus Not Scraping Metrics

**Symptom:** Prometheus shows "Down" for services, or no metrics available

**Diagnosis:**
```bash
# Access Prometheus UI
# http://localhost:9090/targets

# Check if endpoints are responding
curl http://localhost:8080/metrics
curl http://localhost:8081/metrics
```

**Solutions:**

**1. Service Not Exposing Metrics**
```bash
# Check if service implements /metrics endpoint
curl http://localhost:8080/metrics

# If 404, add metrics handler to service (see DEVELOPMENT_GUIDE.md)
```

**2. Prometheus Configuration Wrong**
```bash
# Check prometheus.yml
cat infra/prometheus/prometheus.yml

# Verify service names and ports match docker-compose.yml
# Restart Prometheus
docker-compose restart prometheus
```

**3. Service Unreachable from Prometheus**
```bash
# Test connectivity from Prometheus container
docker-compose exec prometheus wget -O- http://workout-service:8080/metrics

# If fails, services might not be on same network
# Check docker-compose.yml - all should be on 'ironlog' network
```

### Grafana Dashboards Empty

**Symptom:** Grafana shows no data in panels

**Diagnosis:**
```bash
# Check Prometheus datasource connection
# Grafana UI → Configuration → Data Sources → Prometheus
# Click "Save & Test"

# Query Prometheus directly
# http://localhost:9090/graph
# Try: up{job="workout-service"}
```

**Solutions:**

**1. Metrics Not Being Collected**
```bash
# Check if service is exporting metrics
curl http://localhost:8080/metrics

# Should return Prometheus format:
# # TYPE metric_name counter
# metric_name 0
```

**2. Datasource Not Connected**
```bash
# Delete and recreate datasource
# Grafana UI → Configuration → Data Sources
# Delete Prometheus datasource
# Add new: URL=http://prometheus:9090
```

### Jaeger Not Showing Traces

**Symptom:** Jaeger UI empty or "No traces found"

**Diagnosis:**
```bash
# Check if services are sending traces
docker-compose logs parser-service | grep -i "trace"

# Check Jaeger UI
# http://localhost:16686

# Verify services have tracing enabled
```

**Solutions:**

**1. Tracing Not Enabled in Services**
```bash
# Add tracing initialization to service (see ARCHITECTURE.md)
# Usually in cmd/main.go:

# jaegerExporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(...))
# tp := tracesdk.NewTracerProvider(tracesdk.WithBatcher(jaegerExporter))
```

**2. Jaeger Collector Not Ready**
```bash
# Check Jaeger logs
docker-compose logs jaeger

# Restart
docker-compose restart jaeger
sleep 10
```

---

## Frontend Issues

### Frontend Cannot Connect to Backend

**Symptom:** "Failed to fetch" in browser console

**Diagnosis:**
```bash
# Check browser console (F12)
# Look for CORS errors or connection refused

# Test backend directly
curl http://localhost:8080/health

# Test from frontend container
docker-compose exec web-app curl http://workout-service:8080/health
```

**Solutions:**

**1. Backend Not Running**
```bash
docker-compose start workout-service
```

**2. CORS Configuration Missing**
```bash
# Add CORS headers in service
# In cmd/main.go or middleware:

w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

**3. Wrong API URL in Frontend**
```bash
# Check frontend/web-app/src/services/api.ts
# Should point to: http://localhost:8080 (for local)
# Or correct backend service in docker network

# Example:
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
```

### Frontend Build Fails

**Symptom:** `npm run build` fails with errors

**Diagnosis:**
```bash
# Check Node version
node --version  # Should be 18+

# Check npm install
npm --version

# Try clean build
rm -rf node_modules package-lock.json
npm install
npm run build
```

**Solutions:**

**1. TypeScript Errors**
```bash
# Check for type errors
npx tsc --noEmit

# Fix errors or add type definitions
```

**2. Missing Dependencies**
```bash
npm install
npm run build
```

**3. Build Environment Variables**
```bash
# Create .env file in frontend/web-app/
VITE_API_URL=http://localhost:8080

# Use in code:
const apiUrl = import.meta.env.VITE_API_URL;
```

---

## Full System Reset

When all else fails:

```bash
# Stop everything
docker-compose down

# Remove all volumes (data loss!)
docker volume rm $(docker volume ls -q | grep ironlog)

# Remove images
docker-compose down --rmi all

# Start fresh
docker-compose up -d

# Wait for initialization
sleep 60

# Verify
docker-compose ps  # All should be Up
curl http://localhost:8080/health
```

---

## Getting Help

### Collect Diagnostic Information

```bash
# System info
docker --version
docker-compose --version

# Container status
docker-compose ps

# All logs (last 100 lines)
docker-compose logs --tail 100 > diagnostic.log

# Resource usage
docker stats --no-stream

# Database info
docker-compose exec postgres psql -U postgres -d ironlog -c "\d"

# Kafka info
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:29092
```

### Check Documentation

1. [QUICKSTART.md](QUICKSTART.md) - Setup issues
2. [DEPLOYMENT_DOCKER.md](DEPLOYMENT_DOCKER.md) - Docker configuration
3. [SERVICE_COMMUNICATION.md](SERVICE_COMMUNICATION.md) - Event flow issues
4. [E2E_TEST_GUIDE.md](E2E_TEST_GUIDE.md) - Step-by-step verification

### Debug Checklist

- [ ] All containers running (`docker-compose ps`)
- [ ] PostgreSQL accessible (`psql` connection)
- [ ] Kafka topics created (`kafka-topics --list`)
- [ ] Service endpoints responding (`curl /health`)
- [ ] Metrics being collected (`curl /metrics`)
- [ ] Events flowing through system (`kafka-console-consumer`)
- [ ] Frontend can reach backend (browser network tab)
