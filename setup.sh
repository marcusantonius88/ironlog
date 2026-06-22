#!/bin/bash

# IRONLOG Project Setup Script
# Initializes the complete project development environment with all prerequisites,
# dependencies, and Docker infrastructure.

set -e

echo "🏋️  IRONLOG - Setting up development environment..."
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose."
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+."
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+."
    exit 1
fi

echo "✅ All prerequisites found"
echo ""

# Create environment files
echo "🔧 Setting up environment files..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file"
else
    echo "✅ .env already exists"
fi

# Initialize frontend
echo "📦 Setting up frontend dependencies..."
cd frontend/web-app
npm install
cd ../..
echo "✅ Frontend setup complete"

# Start infrastructure
echo "🐳 Starting Docker infrastructure..."
cd infra
docker-compose up -d
cd ..
echo "✅ Infrastructure started"

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

echo ""
echo "🎉 Setup complete!"
echo ""
echo "📍 Services available at:"
echo "   Frontend:    http://localhost:3000"
echo "   Workout API: http://localhost:8080"
echo "   Parser API:  http://localhost:8081"
echo "   Prometheus:  http://localhost:9090"
echo "   Grafana:     http://localhost:3001 (admin/password)"
echo "   Jaeger:      http://localhost:16686"
echo ""
echo "Run 'make docker-down' to stop services"
echo ""
