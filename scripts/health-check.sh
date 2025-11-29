#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🏥 LLM Observability Platform - Health Check"
echo "=============================================="
echo ""

# Function to check if a service is healthy
check_service() {
    local service_name=$1
    local check_command=$2
    
    echo -n "Checking $service_name... "
    
    if eval "$check_command" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Healthy${NC}"
        return 0
    else
        echo -e "${RED}✗ Unhealthy${NC}"
        return 1
    fi
}

# Check Docker
echo "📦 Docker Services:"
check_service "Docker" "docker ps"

# Check ClickHouse
echo ""
echo "💾 Databases:"
check_service "ClickHouse HTTP" "curl -s http://localhost:8123"
check_service "ClickHouse Query" "curl -s 'http://localhost:8123/?query=SELECT%201'"

# Check Redis (FIXED - using Docker)
check_service "Redis" "docker exec llm-obs-redis redis-cli ping"

# Check Kafka
echo ""
echo "📨 Message Queue:"
check_service "Kafka" "docker exec llm-obs-kafka kafka-broker-api-versions --bootstrap-server localhost:9092"

# Check Prometheus
echo ""
echo "📊 Monitoring:"
check_service "Prometheus" "curl -s http://localhost:9090/-/healthy"

# Check Grafana
check_service "Grafana" "curl -s http://localhost:3001/api/health"

echo ""
echo "=============================================="
echo "Health check complete!"
