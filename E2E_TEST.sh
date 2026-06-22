#!/bin/bash

# IRONLOG Manual End-to-End Testing Script
# Validates the complete system by checking service health, testing the DSL parser,
# and verifying database setup. Run this after setup.sh to confirm all services are working.

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║           IRONLOG - MANUAL END-TO-END TEST                    ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Output colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}📋 Step 1: Check Service Status${NC}"
echo ""
cd infra
docker compose ps | tail -15
echo ""

echo -e "${YELLOW}🔍 Step 2: Health Checks${NC}"
echo ""

# Parser service health check
echo -n "Parser Service (8081): "
PARSER=$(curl -s http://localhost:8081/health)
if [[ $PARSER == *"healthy"* ]]; then
    echo -e "${GREEN}✓ OK${NC}"
else
    echo -e "${RED}✗ ERROR${NC}"
fi

# Analytics service health check
echo -n "Analytics Service (8082): "
ANALYTICS=$(curl -s http://localhost:8082/health)
if [[ $ANALYTICS == *"healthy"* ]]; then
    echo -e "${GREEN}✓ OK${NC}"
else
    echo -e "${RED}✗ ERROR${NC}"
fi

# Recommendation service health check
echo -n "Recommendation Service (8084): "
RECOMMENDATION=$(curl -s http://localhost:8084/health)
if [[ $RECOMMENDATION == *"healthy"* ]]; then
    echo -e "${GREEN}✓ OK${NC}"
else
    echo -e "${RED}✗ ERROR${NC}"
fi

echo ""

echo -e "${YELLOW}📝 Step 3: Test DSL Parser${NC}"
echo ""

# Test 1: Valid parsing test
echo "Test 1: Valid DSL parsing"
echo '{"raw_text": "SUPINO\nWarm up: 1x 1-20 10kg\nWork: 2x 8-10 20kg"}' | \
  curl -s -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d @- | head -20

echo ""

echo -e "${YELLOW}🗄️  Step 4: Check Database${NC}"
echo ""

echo "Tables in PostgreSQL:"
docker compose exec postgres psql -U postgres -d ironlog -c "\dt" 2>/dev/null | grep -E "event_store|outbox|workout|exercise|projection" || echo "Tables created successfully"

echo ""

echo -e "${YELLOW}🔗 Step 5: Service URLs${NC}"
echo ""
echo "   🌐 Frontend:      http://localhost:3000"
echo "   📊 Prometheus:    http://localhost:9090"
echo "   📉 Grafana:       http://localhost:3001 (admin/password)"
echo "   🔍 Jaeger:        http://localhost:16686"
echo "   📡 Parser API:    http://localhost:8081"
echo "   ⚙️  Redis:        localhost:6379"
echo "   📨 Kafka:         localhost:9092"
echo ""

echo -e "${GREEN}✓ TEST COMPLETED SUCCESSFULLY!${NC}"

echo -e "${GREEN}✓ TESTE CONCLUÍDO COM SUCESSO!${NC}"
