#!/bin/bash

# Automated Integration Test Runner
# Starts server, runs tests, cleans up

set -e  # Exit on error

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Automated Integration Test Runner"
echo "===================================="

# Check if server is already running
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Server already running on port 8080${NC}"
    echo "Will use existing server for tests"
    SERVER_WAS_RUNNING=true
else
    echo "📦 Starting server..."
    SERVER_WAS_RUNNING=false
    
    # Start server in background
    go run cmd/api/main.go > /tmp/llm-obs-server.log 2>&1 &
    SERVER_PID=$!
    
    echo "Server PID: $SERVER_PID"
    echo "Waiting for server to be ready..."
    
    # Wait for server to start (max 30 seconds)
    COUNTER=0
    MAX_ATTEMPTS=30
    
    while [ $COUNTER -lt $MAX_ATTEMPTS ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Server is ready!${NC}"
            break
        fi
        echo -n "."
        sleep 1
        ((COUNTER++))
    done
    
    if [ $COUNTER -eq $MAX_ATTEMPTS ]; then
        echo -e "\n${RED}❌ Server failed to start within 30 seconds${NC}"
        echo "Check logs at: /tmp/llm-obs-server.log"
        cat /tmp/llm-obs-server.log
        exit 1
    fi
fi

echo ""
echo "🧪 Running integration tests..."
echo "--------------------------------"

# Run tests
if go test ./tests/integration -v; then
    echo ""
    echo -e "${GREEN}✅ All integration tests passed!${NC}"
    TEST_RESULT=0
else
    echo ""
    echo -e "${RED}❌ Some tests failed${NC}"
    TEST_RESULT=1
fi

# Cleanup - only if we started the server
if [ "$SERVER_WAS_RUNNING" = false ]; then
    echo ""
    echo "🧹 Cleaning up..."
    
    # Kill server
    if [ ! -z "$SERVER_PID" ]; then
        echo "Stopping server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
        
        # Wait for it to stop
        sleep 2
        
        # Force kill if still running
        kill -9 $SERVER_PID 2>/dev/null || true
    fi
    
    echo "✅ Cleanup complete"
fi

echo ""
echo "===================================="
exit $TEST_RESULT