package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sahared/llm-observability/internal/models"
)

// SaveMetric saves a metric to ClickHouse
func (r *ClickHouseRepository) SaveMetric(ctx context.Context, metric *models.Metric) error {
	tagsJSON, err := json.Marshal(metric.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO llm_observability.metrics (
			timestamp, organization_id, project_id,
			metric_name, metric_value, tags
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		metric.Timestamp,
		metric.OrganizationID,
		metric.ProjectID,
		metric.MetricName,
		metric.MetricValue,
		string(tagsJSON),
	)

	if err != nil {
		return fmt.Errorf("failed to insert metric: %w", err)
	}

	return nil
}

// GetMetrics retrieves metrics based on query parameters
func (r *ClickHouseRepository) GetMetrics(ctx context.Context, query *models.MetricQuery) ([]*models.Metric, error) {
	sqlQuery := `
		SELECT 
			timestamp, organization_id, project_id,
			metric_name, metric_value, tags
		FROM llm_observability.metrics
		WHERE organization_id = ?
		  AND project_id = ?
		  AND metric_name = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
		ORDER BY timestamp DESC
		LIMIT 1000
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery,
		query.OrganizationID,
		query.ProjectID,
		query.MetricName,
		query.StartTime,
		query.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*models.Metric
	for rows.Next() {
		var metric models.Metric
		var tagsJSON string

		err := rows.Scan(
			&metric.Timestamp,
			&metric.OrganizationID,
			&metric.ProjectID,
			&metric.MetricName,
			&metric.MetricValue,
			&tagsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &metric.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// GetMetricSummary retrieves aggregated metrics for a time range
func (r *ClickHouseRepository) GetMetricSummary(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*models.MetricSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_requests,
			AVG(duration_ms) as avg_latency,
			quantile(0.50)(duration_ms) as p50_latency,
			quantile(0.95)(duration_ms) as p95_latency,
			quantile(0.99)(duration_ms) as p99_latency,
			SUM(total_cost_usd) as total_cost,
			AVG(total_cost_usd) as avg_cost,
			SUM(total_tokens) as total_tokens,
			countIf(status = 'success') * 100.0 / COUNT(*) as success_rate,
			countIf(status = 'error') * 100.0 / COUNT(*) as error_rate
		FROM llm_observability.traces
		WHERE organization_id = ?
		  AND project_id = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
	`

	var summary models.MetricSummary
	err := r.db.QueryRowContext(ctx, query, orgID, projectID, startTime, endTime).Scan(
		&summary.TotalRequests,
		&summary.AvgLatencyMs,
		&summary.P50LatencyMs,
		&summary.P95LatencyMs,
		&summary.P99LatencyMs,
		&summary.TotalCostUSD,
		&summary.AvgCostPerRequest,
		&summary.TotalTokens,
		&summary.SuccessRate,
		&summary.ErrorRate,
	)

	if err == sql.ErrNoRows {
		// No data found, return empty summary
		return &models.MetricSummary{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query metric summary: %w", err)
	}

	return &summary, nil
}

// GetCostBreakdown retrieves cost breakdown by model
func (r *ClickHouseRepository) GetCostBreakdown(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.CostBreakdown, error) {
	query := `
		SELECT 
			model,
			SUM(total_cost_usd) as total_cost,
			COUNT(*) as total_calls,
			AVG(total_cost_usd) as avg_cost
		FROM llm_observability.traces
		WHERE organization_id = ?
		  AND project_id = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
		GROUP BY model
		ORDER BY total_cost DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query cost breakdown: %w", err)
	}
	defer rows.Close()

	var breakdown []*models.CostBreakdown
	var totalCost float64

	// First pass: collect data and calculate total
	for rows.Next() {
		var cb models.CostBreakdown
		err := rows.Scan(&cb.Model, &cb.TotalCost, &cb.TotalCalls, &cb.AvgCost)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cost breakdown: %w", err)
		}
		totalCost += cb.TotalCost
		breakdown = append(breakdown, &cb)
	}

	// Second pass: calculate percentages
	for _, cb := range breakdown {
		if totalCost > 0 {
			cb.Percentage = (cb.TotalCost / totalCost) * 100.0
		}
	}

	return breakdown, nil
}

// GetModelUsage retrieves usage statistics per model
func (r *ClickHouseRepository) GetModelUsage(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.ModelUsage, error) {
	query := `
		SELECT 
			model,
			COUNT(*) as call_count,
			SUM(total_tokens) as total_tokens,
			SUM(total_cost_usd) as total_cost,
			AVG(duration_ms) as avg_latency
		FROM llm_observability.traces
		WHERE organization_id = ?
		  AND project_id = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
		GROUP BY model
		ORDER BY call_count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query model usage: %w", err)
	}
	defer rows.Close()

	var usage []*models.ModelUsage
	for rows.Next() {
		var mu models.ModelUsage
		err := rows.Scan(&mu.Model, &mu.CallCount, &mu.TotalTokens, &mu.TotalCost, &mu.AvgLatency)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model usage: %w", err)
		}
		usage = append(usage, &mu)
	}

	return usage, nil
}
