#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo "🧪 Testing LLM Observability Platform API"
echo "=========================================="

# Test 1: Health Check
echo ""
echo -e "${YELLOW}Test 1: Health Check${NC}"
RESPONSE=$(curl -s "$BASE_URL/health")
if echo "$RESPONSE" | grep -q "ok"; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo -e "${RED}✗ Health check failed${NC}"
    echo "$RESPONSE"
fi

# Test 2: API Info
echo ""
echo -e "${YELLOW}Test 2: API Info${NC}"
RESPONSE=$(curl -s "$BASE_URL/api/v1")
if echo "$RESPONSE" | grep -q "LLM Observability"; then
    echo -e "${GREEN}✓ API info passed${NC}"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo -e "${RED}✗ API info failed${NC}"
    echo "$RESPONSE"
fi

# Test 3: Create Trace
echo ""
echo -e "${YELLOW}Test 3: Create Trace${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/traces" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "org-test",
    "project_id": "proj-test",
    "trace_type": "single_call",
    "model": "gpt-4",
    "provider": "openai",
    "spans": [{
      "name": "test_span",
      "model": "gpt-4",
      "provider": "openai",
      "input": "What is 2+2?",
      "output": "4",
      "prompt_tokens": 10,
      "completion_tokens": 5,
      "duration_ms": 150,
      "status": "success"
    }]
  }')

if echo "$RESPONSE" | grep -q "accepted"; then
    echo -e "${GREEN}✓ Create trace passed${NC}"
    TRACE_ID=$(echo "$RESPONSE" | jq -r '.data.trace_id' 2>/dev/null)
    echo "  Trace ID: $TRACE_ID"
else
    echo -e "${RED}✗ Create trace failed${NC}"
    echo "$RESPONSE"
fi

# Test 4: List Traces
echo ""
echo -e "${YELLOW}Test 4: List Traces${NC}"
RESPONSE=$(curl -s "$BASE_URL/api/v1/traces?organization_id=org-test&project_id=proj-test&limit=10")
if echo "$RESPONSE" | grep -q "data"; then
    echo -e "${GREEN}✓ List traces passed${NC}"
    COUNT=$(echo "$RESPONSE" | jq -r '.total_count' 2>/dev/null || echo "unknown")
    echo "  Total traces: $COUNT"
else
    echo -e "${RED}✗ List traces failed${NC}"
    echo "$RESPONSE"
fi

# Test 5: Dashboard
echo ""
echo -e "${YELLOW}Test 5: Dashboard Summary${NC}"
RESPONSE=$(curl -s "$BASE_URL/api/v1/analytics/dashboard?organization_id=org-test&project_id=proj-test&time_range=24h")
if echo "$RESPONSE" | grep -q "metric_summary"; then
    echo -e "${GREEN}✓ Dashboard passed${NC}"
    REQUESTS=$(echo "$RESPONSE" | jq -r '.data.metric_summary.total_requests' 2>/dev/null || echo "0")
    echo "  Total requests: $REQUESTS"
else
    echo -e "${RED}✗ Dashboard failed${NC}"
    echo "$RESPONSE"
fi

echo ""
echo "=========================================="
echo "✓ API testing complete!"
