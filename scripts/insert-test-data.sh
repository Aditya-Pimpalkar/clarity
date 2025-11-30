#!/bin/bash

echo "Inserting test data..."

docker exec llm-obs-clickhouse clickhouse-client --query "INSERT INTO llm_observability.organizations VALUES ('org-test-123', 'Test Organization', 'free', now(), now())"
echo "✓ Organization inserted"

docker exec llm-obs-clickhouse clickhouse-client --query "INSERT INTO llm_observability.projects VALUES ('proj-test-456', 'org-test-123', 'Test Project', 'My first observability project', now(), now())"
echo "✓ Project inserted"

docker exec llm-obs-clickhouse clickhouse-client --query "INSERT INTO llm_observability.traces VALUES ('trace-001', 'org-test-123', 'proj-test-456', now(), 'single_call', 245, 'success', 0.0032, 650, 'gpt-4', 'openai', 'user-001', '{\"temperature\": 0.7}')"
echo "✓ Trace inserted"

docker exec llm-obs-clickhouse clickhouse-client --query "INSERT INTO llm_observability.spans VALUES ('span-001', 'trace-001', '', 'llm_completion', now() - INTERVAL 2 SECOND, now(), 245, 'gpt-4', 'openai', 'What is machine learning?', 'Machine learning is AI...', 15, 85, 100, 0.0032, 'success', '', '{}')"
echo "✓ Span inserted"

echo "✅ All test data inserted successfully!"
