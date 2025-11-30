CREATE DATABASE IF NOT EXISTS llm_observability;

USE llm_observability;

CREATE TABLE IF NOT EXISTS traces (
    trace_id String,
    organization_id String,
    project_id String,
    timestamp DateTime64(3),
    trace_type String,
    duration_ms UInt32,
    status String,
    total_cost_usd Float64,
    total_tokens UInt32,
    model String,
    provider String,
    user_id String,
    metadata String,
    INDEX idx_org_project (organization_id, project_id) TYPE minmax GRANULARITY 1,
    INDEX idx_timestamp timestamp TYPE minmax GRANULARITY 1,
    INDEX idx_status status TYPE set(10) GRANULARITY 1,
    INDEX idx_model model TYPE set(100) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (organization_id, project_id, timestamp)
TTL toDateTime(timestamp) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS spans (
    span_id String,
    trace_id String,
    parent_span_id String,
    name String,
    start_time DateTime64(3),
    end_time DateTime64(3),
    duration_ms UInt32,
    model String,
    provider String,
    input String,
    output String,
    prompt_tokens UInt32,
    completion_tokens UInt32,
    total_tokens UInt32,
    cost_usd Float64,
    status String,
    error_message String,
    metadata String,
    INDEX idx_trace_id trace_id TYPE minmax GRANULARITY 1,
    INDEX idx_model model TYPE set(100) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(start_time)
ORDER BY (trace_id, start_time)
TTL toDateTime(start_time) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS metrics (
    timestamp DateTime64(3),
    organization_id String,
    project_id String,
    metric_name String,
    metric_value Float64,
    tags String,
    INDEX idx_org_project (organization_id, project_id) TYPE minmax GRANULARITY 1,
    INDEX idx_metric_name metric_name TYPE set(50) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (organization_id, project_id, metric_name, timestamp)
TTL toDateTime(timestamp) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS users (
    id String,
    organization_id String,
    email String,
    name String,
    role String,
    created_at DateTime64(3),
    updated_at DateTime64(3),
    INDEX idx_org organization_id TYPE minmax GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY (organization_id, email)
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS organizations (
    id String,
    name String,
    plan String,
    created_at DateTime64(3),
    updated_at DateTime64(3)
) ENGINE = MergeTree()
ORDER BY id
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS projects (
    id String,
    organization_id String,
    name String,
    description String,
    created_at DateTime64(3),
    updated_at DateTime64(3),
    INDEX idx_org organization_id TYPE minmax GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY (organization_id, id)
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS api_keys (
    id String,
    organization_id String,
    project_id String,
    key_hash String,
    name String,
    created_at DateTime64(3),
    last_used_at Nullable(DateTime64(3)),
    INDEX idx_org_project (organization_id, project_id) TYPE minmax GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY (organization_id, project_id, id)
SETTINGS index_granularity = 8192;

CREATE MATERIALIZED VIEW IF NOT EXISTS metrics_hourly
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(hour)
ORDER BY (organization_id, project_id, metric_name, hour)
AS SELECT
    toStartOfHour(timestamp) AS hour,
    organization_id,
    project_id,
    metric_name,
    sum(metric_value) AS total_value,
    count() AS count
FROM metrics
GROUP BY hour, organization_id, project_id, metric_name;

CREATE MATERIALIZED VIEW IF NOT EXISTS daily_costs
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(day)
ORDER BY (organization_id, project_id, day)
AS SELECT
    toDate(timestamp) AS day,
    organization_id,
    project_id,
    sum(total_cost_usd) AS total_cost,
    count() AS trace_count
FROM traces
GROUP BY day, organization_id, project_id;
