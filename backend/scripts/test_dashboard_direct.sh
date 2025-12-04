#!/bin/bash

echo "Testing dashboard with direct database query..."
echo ""

# Get organization IDs
echo "Organizations in database:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT id, name FROM organizations FORMAT Pretty"

echo ""
echo "Testing with org-test-123..."
echo ""

# Test direct SQL queries that the dashboard would use
echo "1. Total Traces:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT count() FROM traces WHERE organization_id = 'org-test-123'"

echo ""
echo "2. Total Cost:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT sum(total_cost_usd) FROM traces WHERE organization_id = 'org-test-123'"

echo ""
echo "3. Avg Latency:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT avg(duration_ms) FROM traces WHERE organization_id = 'org-test-123'"

echo ""
echo "4. Success Rate:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT countIf(status = 'success') * 100.0 / count() FROM traces WHERE organization_id = 'org-test-123'"

echo ""
echo "5. Top Models:"
docker exec -i llm-obs-clickhouse clickhouse-client --database=llm_observability \
  --query "SELECT model, count() as cnt, sum(total_cost_usd) as cost FROM traces WHERE organization_id = 'org-test-123' GROUP BY model ORDER BY cnt DESC FORMAT Pretty"
