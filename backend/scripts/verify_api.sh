#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0

echo "🧪 API Verification Test Suite"
echo "=============================="
echo ""

# Test function
test_endpoint() {
    local name=$1
    local url=$2
    local expected=$3
    
    echo -n "Testing $name... "
    
    response=$(curl -s "$url")
    
    if echo "$response" | grep -q "$expected"; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "  Response: $response"
        ((FAILED++))
    fi
}

# Health checks
echo "1️⃣  Health Checks"
echo "─────────────────"
test_endpoint "Health" "$BASE_URL/health" '"status":"ok"'
test_endpoint "Readiness" "$BASE_URL/ready" '"ready":true'
test_endpoint "Liveness" "$BASE_URL/live" '"alive":true'
echo ""

# API info
echo "2️⃣  API Info"
echo "─────────────────"
test_endpoint "API Info" "$BASE_URL/api/v1" "LLM Observability"
echo ""

# Create trace
echo "3️⃣  Trace Creation"
echo "─────────────────"
ORG_ID="org-verify-$(date +%s)"
PROJ_ID="proj-verify-$(date +%s)"

echo -n "Creating test trace... "
response=$(curl -s -X POST "$BASE_URL/api/v1/traces" \
  -H "Content-Type: application/json" \
  -d "{
    \"organization_id\": \"$ORG_ID\",
    \"project_id\": \"$PROJ_ID\",
    \"trace_type\": \"single_call\",
    \"model\": \"gpt-4\",
    \"provider\": \"openai\",
    \"spans\": [{
      \"name\": \"test\",
      \"model\": \"gpt-4\",
      \"provider\": \"openai\",
      \"input\": \"test\",
      \"output\": \"test\",
      \"prompt_tokens\": 10,
      \"completion_tokens\": 5,
      \"duration_ms\": 100,
      \"status\": \"success\"
    }]
  }")

if echo "$response" | grep -q '"status":"accepted"'; then
    echo -e "${GREEN}✓ PASSED${NC}"
    ((PASSED++))
    TRACE_ID=$(echo "$response" | grep -o '"trace_id":"[^"]*"' | cut -d'"' -f4)
    echo "  Trace ID: $TRACE_ID"
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "  Response: $response"
    ((FAILED++))
fi
echo ""

# Get trace
echo "4️⃣  Trace Retrieval"
echo "─────────────────"
if [ -n "$TRACE_ID" ]; then
    test_endpoint "Get Trace" "$BASE_URL/api/v1/traces/$TRACE_ID" '"trace_id":"'
else
    echo -e "${YELLOW}⊘ SKIPPED (no trace ID)${NC}"
fi

test_endpoint "List Traces" "$BASE_URL/api/v1/traces?organization_id=$ORG_ID&project_id=$PROJ_ID" '"data":'
echo ""

# Analytics
echo "5️⃣  Analytics"
echo "─────────────────"
test_endpoint "Dashboard" "$BASE_URL/api/v1/analytics/dashboard?organization_id=$ORG_ID&project_id=$PROJ_ID&time_range=24h" "metric_summary"
test_endpoint "Costs" "$BASE_URL/api/v1/analytics/costs?organization_id=$ORG_ID&project_id=$PROJ_ID" "cost_breakdown"
test_endpoint "Performance" "$BASE_URL/api/v1/analytics/performance?organization_id=$ORG_ID&project_id=$PROJ_ID" '"summary"'
test_endpoint "Models" "$BASE_URL/api/v1/analytics/models?organization_id=$ORG_ID&project_id=$PROJ_ID" '"model"'
echo ""

# Summary
echo "=============================="
echo "📊 Test Summary"
echo "=============================="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
else
    echo -e "${GREEN}Failed: $FAILED${NC}"
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed! Your API is working perfectly!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Check the output above for details.${NC}"
    exit 1
fi
