package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sahared/llm-observability/internal/models"
)

// ClickHouseRepository implements Repository interface for ClickHouse
type ClickHouseRepository struct {
	db *sql.DB
}

// formatTimeForClickHouse converts time to ClickHouse compatible format
func formatTimeForClickHouse(timeStr string) string {
	// Try to parse as RFC3339
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// If parsing fails, return as-is
		return timeStr
	}
	// Convert to UTC and format for ClickHouse (without timezone)
	return t.UTC().Format("2006-01-02 15:04:05")
}

// NewClickHouseRepository creates a new ClickHouse repository
func NewClickHouseRepository(dsn string) (*ClickHouseRepository, error) {
	// Open connection to ClickHouse
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &ClickHouseRepository{db: db}, nil
}

// Ping checks if the database connection is alive
func (r *ClickHouseRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Close closes the database connection
func (r *ClickHouseRepository) Close() error {
	return r.db.Close()
}

// SaveTrace saves a trace and all its spans to ClickHouse
func (r *ClickHouseRepository) SaveTrace(ctx context.Context, trace *models.Trace) error {
	// Start a transaction (ClickHouse doesn't support full ACID transactions,
	// but we use this for batch inserts)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we don't commit

	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(trace.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Insert trace
	query := `
		INSERT INTO llm_observability.traces (
			trace_id, organization_id, project_id, timestamp,
			trace_type, duration_ms, status, total_cost_usd,
			total_tokens, model, provider, user_id, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, query,
		trace.TraceID,
		trace.OrganizationID,
		trace.ProjectID,
		trace.Timestamp,
		trace.TraceType,
		trace.DurationMs,
		trace.Status,
		trace.TotalCostUSD,
		trace.TotalTokens,
		trace.Model,
		trace.Provider,
		trace.UserID,
		string(metadataJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert trace: %w", err)
	}

	// Insert all spans
	for _, span := range trace.Spans {
		if err := r.saveSpanInTx(ctx, tx, &span); err != nil {
			return fmt.Errorf("failed to insert span %s: %w", span.SpanID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// saveSpanInTx saves a span within a transaction
func (r *ClickHouseRepository) saveSpanInTx(ctx context.Context, tx *sql.Tx, span *models.Span) error {
	metadataJSON, err := json.Marshal(span.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal span metadata: %w", err)
	}

	query := `
		INSERT INTO llm_observability.spans (
			span_id, trace_id, parent_span_id, name,
			start_time, end_time, duration_ms, model, provider,
			input, output, prompt_tokens, completion_tokens,
			total_tokens, cost_usd, status, error_message, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, query,
		span.SpanID,
		span.TraceID,
		span.ParentSpanID,
		span.Name,
		span.StartTime,
		span.EndTime,
		span.DurationMs,
		span.Model,
		span.Provider,
		span.Input,
		span.Output,
		span.PromptTokens,
		span.CompletionTokens,
		span.TotalTokens,
		span.CostUSD,
		span.Status,
		span.ErrorMessage,
		string(metadataJSON),
	)

	return err
}

// SaveSpan saves a single span (public method)
func (r *ClickHouseRepository) SaveSpan(ctx context.Context, span *models.Span) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := r.saveSpanInTx(ctx, tx, span); err != nil {
		return err
	}

	return tx.Commit()
}

// GetTraceByID retrieves a single trace by its ID
func (r *ClickHouseRepository) GetTraceByID(ctx context.Context, traceID string) (*models.Trace, error) {
	query := `
		SELECT 
			trace_id, organization_id, project_id, timestamp,
			trace_type, duration_ms, status, total_cost_usd,
			total_tokens, model, provider, user_id, metadata
		FROM llm_observability.traces
		WHERE trace_id = ?
		LIMIT 1
	`

	var trace models.Trace
	var metadataJSON string

	err := r.db.QueryRowContext(ctx, query, traceID).Scan(
		&trace.TraceID,
		&trace.OrganizationID,
		&trace.ProjectID,
		&trace.Timestamp,
		&trace.TraceType,
		&trace.DurationMs,
		&trace.Status,
		&trace.TotalCostUSD,
		&trace.TotalTokens,
		&trace.Model,
		&trace.Provider,
		&trace.UserID,
		&metadataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("trace not found: %s", traceID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query trace: %w", err)
	}

	// Parse metadata JSON
	if err := json.Unmarshal([]byte(metadataJSON), &trace.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Get spans for this trace
	spans, err := r.GetSpansByTraceID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spans: %w", err)
	}
	trace.Spans = spans

	return &trace, nil
}

// GetSpansByTraceID retrieves all spans for a trace
func (r *ClickHouseRepository) GetSpansByTraceID(ctx context.Context, traceID string) ([]models.Span, error) {
	query := `
		SELECT 
			span_id, trace_id, parent_span_id, name,
			start_time, end_time, duration_ms, model, provider,
			input, output, prompt_tokens, completion_tokens,
			total_tokens, cost_usd, status, error_message, metadata
		FROM llm_observability.spans
		WHERE trace_id = ?
		ORDER BY start_time
	`

	rows, err := r.db.QueryContext(ctx, query, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query spans: %w", err)
	}
	defer rows.Close()

	var spans []models.Span
	for rows.Next() {
		var span models.Span
		var metadataJSON string

		err := rows.Scan(
			&span.SpanID,
			&span.TraceID,
			&span.ParentSpanID,
			&span.Name,
			&span.StartTime,
			&span.EndTime,
			&span.DurationMs,
			&span.Model,
			&span.Provider,
			&span.Input,
			&span.Output,
			&span.PromptTokens,
			&span.CompletionTokens,
			&span.TotalTokens,
			&span.CostUSD,
			&span.Status,
			&span.ErrorMessage,
			&metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan span: %w", err)
		}

		// Parse metadata
		if err := json.Unmarshal([]byte(metadataJSON), &span.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal span metadata: %w", err)
		}

		spans = append(spans, span)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating spans: %w", err)
	}

	return spans, nil
}

// GetTraces retrieves traces with filtering and pagination
func (r *ClickHouseRepository) GetTraces(ctx context.Context, query *models.TraceQuery) ([]*models.Trace, error) {
	// Build the WHERE clause dynamically
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.OrganizationID != "" {
		whereClause += " AND organization_id = ?"
		args = append(args, query.OrganizationID)
	}

	if query.ProjectID != "" {
		whereClause += " AND project_id = ?"
		args = append(args, query.ProjectID)
	}

	if query.StartTime != "" {
		whereClause += " AND timestamp >= ?"
		args = append(args, formatTimeForClickHouse(query.StartTime))
	}

	if query.EndTime != "" {
		whereClause += " AND timestamp <= ?"
		args = append(args, formatTimeForClickHouse(query.EndTime))
	}

	if query.Model != "" {
		whereClause += " AND model = ?"
		args = append(args, query.Model)
	}

	if query.Status != "" {
		whereClause += " AND status = ?"
		args = append(args, query.Status)
	}

	// Set default limit if not specified
	limit := query.Limit
	if limit == 0 {
		limit = 100
	}

	// Build final query
	sqlQuery := fmt.Sprintf(`
		SELECT 
			trace_id, organization_id, project_id, timestamp,
			trace_type, duration_ms, status, total_cost_usd,
			total_tokens, model, provider, user_id, metadata
		FROM llm_observability.traces
		%s
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, query.Offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query traces: %w", err)
	}
	defer rows.Close()

	// Parse results
	var traces []*models.Trace
	for rows.Next() {
		var trace models.Trace
		var metadataJSON string

		err := rows.Scan(
			&trace.TraceID,
			&trace.OrganizationID,
			&trace.ProjectID,
			&trace.Timestamp,
			&trace.TraceType,
			&trace.DurationMs,
			&trace.Status,
			&trace.TotalCostUSD,
			&trace.TotalTokens,
			&trace.Model,
			&trace.Provider,
			&trace.UserID,
			&metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trace: %w", err)
		}

		// Parse metadata
		if err := json.Unmarshal([]byte(metadataJSON), &trace.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		traces = append(traces, &trace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating traces: %w", err)
	}

	return traces, nil
}

// GetTraceCount returns the total count of traces matching the query
func (r *ClickHouseRepository) GetTraceCount(ctx context.Context, query *models.TraceQuery) (int64, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.OrganizationID != "" {
		whereClause += " AND organization_id = ?"
		args = append(args, query.OrganizationID)
	}

	if query.ProjectID != "" {
		whereClause += " AND project_id = ?"
		args = append(args, query.ProjectID)
	}

	if query.StartTime != "" {
		whereClause += " AND timestamp >= ?"
		args = append(args, formatTimeForClickHouse(query.StartTime))
	}

	if query.EndTime != "" {
		whereClause += " AND timestamp <= ?"
		args = append(args, formatTimeForClickHouse(query.EndTime))
	}

	sqlQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM llm_observability.traces
		%s
	`, whereClause)

	var count int64
	err := r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count traces: %w", err)
	}

	return count, nil
}
