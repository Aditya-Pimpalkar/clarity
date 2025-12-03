#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "🧪 LLM Observability Platform - Complete Test Suite"
echo "===================================================="
echo ""

# Check if server is running
echo -e "${YELLOW}Checking if server is running...${NC}"
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Server is running${NC}"
    SERVER_WAS_RUNNING=true
else
    echo -e "${YELLOW}⚠️  Server not running, starting it...${NC}"
    go run cmd/api/main.go > /tmp/test-server.log 2>&1 &
    SERVER_PID=$!
    SERVER_WAS_RUNNING=false
    
    # Wait for server
    echo -n "Waiting for server to start"
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo ""
            echo -e "${GREEN}✅ Server started (PID: $SERVER_PID)${NC}"
            break
        fi
        echo -n "."
        sleep 1
    done
fi

echo ""
echo "=================================================="
echo -e "${YELLOW}1. Running Unit Tests${NC}"
echo "=================================================="
go test ./internal/... -v -race -coverprofile=coverage.out

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Unit tests passed${NC}"
    echo -e "${YELLOW}Coverage:${NC}"
    go tool cover -func=coverage.out | grep total:
else
    echo -e "${RED}❌ Unit tests failed${NC}"
    exit 1
fi

echo ""
echo "=================================================="
echo -e "${YELLOW}2. Running Integration Tests${NC}"
echo "=================================================="
go test ./tests/integration -v

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Integration tests passed${NC}"
else
    echo -e "${RED}❌ Integration tests failed${NC}"
    exit 1
fi

echo ""
echo "=================================================="
echo -e "${YELLOW}3. Running Benchmark Tests${NC}"
echo "=================================================="
go test ./tests/integration -bench=. -benchmem -run=^$ -benchtime=5s

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Benchmark tests passed${NC}"
else
    echo -e "${RED}❌ Benchmark tests failed${NC}"
fi

# Cleanup
if [ "$SERVER_WAS_RUNNING" = false ]; then
    echo ""
    echo -e "${YELLOW}🛑 Stopping test server...${NC}"
    kill $SERVER_PID 2>/dev/null || true
    echo -e "${GREEN}✅ Server stopped${NC}"
fi

echo ""
echo "===================================================="
echo -e "${GREEN}✨ All tests completed!${NC}"
echo "===================================================="
