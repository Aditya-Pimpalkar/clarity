#!/bin/bash

# Load Testing Script
# Requires: apache bench (ab)

BASE_URL="http://localhost:8080"
API_KEY="test-key-123"

echo "🔥 Load Testing LLM Observability Platform"
echo "==========================================="

# Check if ab is installed
if ! command -v ab &> /dev/null; then
    echo "❌ Apache Bench (ab) not found. Install with:"
    echo "   macOS: brew install httpd"
    echo "   Ubuntu: sudo apt-get install apache2-utils"
    exit 1
fi

# Test 1: Health endpoint
echo ""
echo "Test 1: Health Endpoint (1000 requests, 10 concurrent)"
echo "-------------------------------------------------------"
ab -n 1000 -c 10 -q "$BASE_URL/health" | grep -E "Requests per second|Time per request|Failed requests"

# Test 2: API info
echo ""
echo "Test 2: API Info (1000 requests, 20 concurrent)"
echo "-------------------------------------------------------"
ab -n 1000 -c 20 -q "$BASE_URL/api/v1" | grep -E "Requests per second|Time per request|Failed requests"

# Test 3: Trace creation (requires wrk or custom tool)
echo ""
echo "Test 3: Trace Creation Performance"
echo "-------------------------------------------------------"
echo "Creating 100 traces sequentially..."

SUCCESS=0
FAILED=0
TOTAL_TIME=0

for i in {1..100}; do
    START=$(date +%s%N)
    
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -X POST "$BASE_URL/api/v1/traces" \
        -d "{
            \"organization_id\": \"org-perf\",
            \"project_id\": \"proj-perf\",
            \"trace_type\": \"single_call\",
            \"model\": \"gpt-4\",
            \"provider\": \"openai\",
            \"spans\": [{
                \"name\": \"perf_test\",
                \"model\": \"gpt-4\",
                \"provider\": \"openai\",
                \"prompt_tokens\": $((10 + i)),
                \"completion_tokens\": $((5 + i)),
                \"duration_ms\": $((100 + i)),
                \"status\": \"success\"
            }]
        }")
    
    END=$(date +%s%N)
    DURATION=$(( (END - START) / 1000000 ))
    TOTAL_TIME=$((TOTAL_TIME + DURATION))
    
    if [ "$STATUS" == "201" ]; then
        ((SUCCESS++))
    else
        ((FAILED++))
    fi
done

AVG_TIME=$((TOTAL_TIME / 100))

echo "Results:"
echo "  Successful: $SUCCESS"
echo "  Failed: $FAILED"
echo "  Average time: ${AVG_TIME}ms"
echo "  Throughput: $(bc <<< "scale=2; 1000 / $AVG_TIME") req/sec"

echo ""
echo "==========================================="
echo "✅ Load testing complete!"
