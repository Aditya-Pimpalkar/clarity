#!/bin/bash

# Smoke Test Suite
# Quick tests to verify deployment health

BASE_URL="${1:-http://localhost:8080}"

echo "🔍 Running Smoke Tests"
echo "Target: $BASE_URL"
echo "======================"

PASSED=0
FAILED=0

# Test function
test_endpoint() {
    local name=$1
    local method=$2
    local path=$3
    local expected_status=$4
    local extra_args=${5:-""}
    
    echo -n "Testing $name... "
    
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
        -X "$method" \
        $extra_args \
        "$BASE_URL$path")
    
    if [ "$STATUS" == "$expected_status" ]; then
        echo "✅ PASS"
        ((PASSED++))
    else
        echo "❌ FAIL (got $STATUS, expected $expected_status)"
        ((FAILED++))
    fi
}

# Run tests
test_endpoint "Health Check" "GET" "/health" "200"
test_endpoint "Readiness Check" "GET" "/ready" "200"
test_endpoint "Liveness Check" "GET" "/live" "200"
test_endpoint "API Info" "GET" "/api/v1" "200"
test_endpoint "Unauthorized Access" "GET" "/api/v1/traces?organization_id=test&project_id=test" "401"

# Test with API key (fixed formatting)
echo -n "Testing API Key Auth... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST \
    -H "X-API-Key: test-key-123" \
    -H "Content-Type: application/json" \
    --max-time 10 \
    "$BASE_URL/api/v1/traces" \
    -d '{
        "organization_id": "org-smoke",
        "project_id": "proj-smoke",
        "trace_type": "single_call",
        "model": "gpt-4",
        "provider": "openai",
        "spans": [{
            "name": "smoke_test",
            "model": "gpt-4",
            "provider": "openai",
            "prompt_tokens": 10,
            "completion_tokens": 5,
            "duration_ms": 100,
            "status": "success"
        }]
    }')

if [ "$STATUS" == "201" ]; then
    echo "✅ PASS"
    ((PASSED++))
else
    echo "❌ FAIL (got $STATUS, expected 201)"
    ((FAILED++))
    
    # Debug info
    echo ""
    echo "Debug: Testing API key endpoint directly..."
    curl -X POST \
        -H "X-API-Key: test-key-123" \
        -H "Content-Type: application/json" \
        "$BASE_URL/api/v1/traces" \
        -d '{
            "organization_id": "org-smoke",
            "project_id": "proj-smoke",
            "trace_type": "single_call",
            "model": "gpt-4",
            "provider": "openai",
            "spans": [{
                "name": "smoke_test",
                "model": "gpt-4",
                "provider": "openai",
                "prompt_tokens": 10,
                "completion_tokens": 5,
                "duration_ms": 100,
                "status": "success"
            }]
        }' | jq '.' 2>/dev/null || echo "Response not JSON"
fi

echo ""
echo "======================"
echo "Results: $PASSED passed, $FAILED failed"

if [ $FAILED -eq 0 ]; then
    echo "✅ All smoke tests passed!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi