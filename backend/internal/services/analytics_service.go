package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sahared/llm-observability/internal/models"
	"github.com/sahared/llm-observability/internal/repository"
)

// AnalyticsService handles business logic for analytics and metrics
type AnalyticsService struct {
	repo repository.Repository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo repository.Repository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// GetDashboardSummary returns a comprehensive dashboard summary
func (s *AnalyticsService) GetDashboardSummary(ctx context.Context, orgID, projectID string, timeRange string) (*DashboardSummary, error) {
	// Validate inputs
	if orgID == "" || projectID == "" {
		return nil, fmt.Errorf("organization_id and project_id are required")
	}

	// Parse time range
	startTime, endTime, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// Get metric summary
	metricSummary, err := s.repo.GetMetricSummary(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric summary: %w", err)
	}

	// Get cost breakdown
	costBreakdown, err := s.repo.GetCostBreakdown(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost breakdown: %w", err)
	}

	// Get model usage
	modelUsage, err := s.repo.GetModelUsage(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get model usage: %w", err)
	}

	// Calculate additional insights
	insights := s.generateInsights(metricSummary, costBreakdown, modelUsage)

	return &DashboardSummary{
		TimeRange:     timeRange,
		StartTime:     startTime,
		EndTime:       endTime,
		MetricSummary: metricSummary,
		CostBreakdown: costBreakdown,
		ModelUsage:    modelUsage,
		Insights:      insights,
	}, nil
}

// parseTimeRange converts time range string to start and end times
func (s *AnalyticsService) parseTimeRange(timeRange string) (time.Time, time.Time, error) {
	now := time.Now()
	var startTime time.Time

	switch timeRange {
	case "1h", "last_hour":
		startTime = now.Add(-1 * time.Hour)
	case "24h", "last_24h", "today":
		startTime = now.Add(-24 * time.Hour)
	case "7d", "last_7d", "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "30d", "last_30d", "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	case "90d", "last_90d":
		startTime = now.Add(-90 * 24 * time.Hour)
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported time range: %s", timeRange)
	}

	return startTime, now, nil
}

// generateInsights generates insights from the data
func (s *AnalyticsService) generateInsights(
	summary *models.MetricSummary,
	costBreakdown []*models.CostBreakdown,
	modelUsage []*models.ModelUsage,
) []Insight {
	insights := []Insight{}

	// Insight 1: Error rate
	if summary.ErrorRate > 5.0 {
		insights = append(insights, Insight{
			Type:        "warning",
			Category:    "reliability",
			Title:       "High Error Rate",
			Description: fmt.Sprintf("Error rate is %.2f%%, which is above the recommended 5%% threshold", summary.ErrorRate),
			Severity:    "high",
		})
	} else if summary.ErrorRate > 1.0 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "reliability",
			Title:       "Elevated Error Rate",
			Description: fmt.Sprintf("Error rate is %.2f%%, consider investigating recent changes", summary.ErrorRate),
			Severity:    "medium",
		})
	}

	// Insight 2: Cost concentration
	if len(costBreakdown) > 0 && costBreakdown[0].Percentage > 80.0 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "cost",
			Title:       "Cost Concentration",
			Description: fmt.Sprintf("%.1f%% of costs come from %s. Consider model optimization or caching", costBreakdown[0].Percentage, costBreakdown[0].Model),
			Severity:    "low",
		})
	}

	// Insight 3: Latency
	if summary.P95LatencyMs > 2000 {
		insights = append(insights, Insight{
			Type:        "warning",
			Category:    "performance",
			Title:       "High P95 Latency",
			Description: fmt.Sprintf("P95 latency is %.0fms. Users may experience slow responses", summary.P95LatencyMs),
			Severity:    "medium",
		})
	}

	// Insight 4: Cost per request
	if summary.AvgCostPerRequest > 0.01 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "cost",
			Title:       "High Average Cost",
			Description: fmt.Sprintf("Average cost per request is $%.4f. Consider optimizing prompts or using cheaper models", summary.AvgCostPerRequest),
			Severity:    "low",
		})
	}

	// Insight 5: Success rate achievement
	if summary.SuccessRate > 99.0 {
		insights = append(insights, Insight{
			Type:        "success",
			Category:    "reliability",
			Title:       "Excellent Reliability",
			Description: fmt.Sprintf("Success rate is %.2f%% - great job!", summary.SuccessRate),
			Severity:    "info",
		})
	}

	// Insight 6: Model usage patterns
	if len(modelUsage) > 0 {
		// Find the most used model
		mostUsed := modelUsage[0]
		for _, usage := range modelUsage {
			if usage.CallCount > mostUsed.CallCount {
				mostUsed = usage
			}
		}

		if mostUsed.CallCount > 100 {
			insights = append(insights, Insight{
				Type:        "info",
				Category:    "usage",
				Title:       "High Model Usage",
				Description: fmt.Sprintf("%s is your most used model with %d calls. Ensure you're getting the best value", mostUsed.Model, mostUsed.CallCount),
				Severity:    "info",
			})
		}
	}

	// Insight 7: Token efficiency
	if len(modelUsage) > 0 {
		for _, usage := range modelUsage {
			if usage.CallCount > 0 {
				avgTokensPerCall := float64(usage.TotalTokens) / float64(usage.CallCount)
				if avgTokensPerCall > 2000 {
					insights = append(insights, Insight{
						Type:        "info",
						Category:    "efficiency",
						Title:       "High Token Usage",
						Description: fmt.Sprintf("%s averages %.0f tokens per call. Consider prompt optimization", usage.Model, avgTokensPerCall),
						Severity:    "low",
					})
					break // Only show one token efficiency insight
				}
			}
		}
	}

	// Insight 8: Model diversity
	if len(modelUsage) == 1 && modelUsage[0].CallCount > 50 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "optimization",
			Title:       "Single Model Usage",
			Description: fmt.Sprintf("You're only using %s. Consider testing other models for cost/performance optimization", modelUsage[0].Model),
			Severity:    "low",
		})
	} else if len(modelUsage) >= 3 {
		insights = append(insights, Insight{
			Type:        "success",
			Category:    "optimization",
			Title:       "Good Model Diversity",
			Description: fmt.Sprintf("Using %d different models - great job optimizing for different use cases!", len(modelUsage)),
			Severity:    "info",
		})
	}

	return insights
}

// GetCostAnalysis returns detailed cost analysis
func (s *AnalyticsService) GetCostAnalysis(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*CostAnalysis, error) {
	// Get cost breakdown by model
	breakdown, err := s.repo.GetCostBreakdown(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost breakdown: %w", err)
	}

	// Calculate total cost
	var totalCost float64
	for _, item := range breakdown {
		totalCost += item.TotalCost
	}

	// Calculate daily average
	days := endTime.Sub(startTime).Hours() / 24
	if days == 0 {
		days = 1
	}
	dailyAverage := totalCost / days

	// Project monthly cost
	monthlyProjection := dailyAverage * 30

	// Find most expensive model
	var mostExpensiveModel string
	var highestCost float64
	if len(breakdown) > 0 {
		mostExpensiveModel = breakdown[0].Model
		highestCost = breakdown[0].TotalCost
	}

	return &CostAnalysis{
		TotalCost:          totalCost,
		DailyAverage:       dailyAverage,
		MonthlyProjection:  monthlyProjection,
		CostBreakdown:      breakdown,
		MostExpensiveModel: mostExpensiveModel,
		HighestCost:        highestCost,
	}, nil
}

// GetPerformanceMetrics returns performance analysis
func (s *AnalyticsService) GetPerformanceMetrics(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*PerformanceMetrics, error) {
	summary, err := s.repo.GetMetricSummary(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric summary: %w", err)
	}

	// Calculate additional metrics
	var status string
	var recommendation string

	// Determine overall performance status
	if summary.P95LatencyMs < 500 && summary.ErrorRate < 1.0 {
		status = "excellent"
		recommendation = "Performance is optimal. Continue monitoring."
	} else if summary.P95LatencyMs < 1000 && summary.ErrorRate < 5.0 {
		status = "good"
		recommendation = "Performance is acceptable but could be improved."
	} else if summary.P95LatencyMs < 2000 && summary.ErrorRate < 10.0 {
		status = "fair"
		recommendation = "Performance issues detected. Consider optimization."
	} else {
		status = "poor"
		recommendation = "Critical performance issues. Immediate attention required."
	}

	return &PerformanceMetrics{
		Summary:        summary,
		Status:         status,
		Recommendation: recommendation,
	}, nil
}

// GetModelComparison compares performance and cost across models
func (s *AnalyticsService) GetModelComparison(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*ModelComparison, error) {
	modelUsage, err := s.repo.GetModelUsage(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get model usage: %w", err)
	}

	comparisons := make([]*ModelComparison, len(modelUsage))
	for i, usage := range modelUsage {
		// Calculate metrics per request
		avgCostPerRequest := 0.0
		avgTokensPerRequest := 0.0
		if usage.CallCount > 0 {
			avgCostPerRequest = usage.TotalCost / float64(usage.CallCount)
			avgTokensPerRequest = float64(usage.TotalTokens) / float64(usage.CallCount)
		}

		// Calculate efficiency score (lower is better)
		// Score = (normalized_latency + normalized_cost) / 2
		efficiencyScore := (usage.AvgLatency / 1000.0) + (avgCostPerRequest * 100.0)

		comparisons[i] = &ModelComparison{
			Model:               usage.Model,
			TotalCalls:          usage.CallCount,
			TotalCost:           usage.TotalCost,
			AvgCostPerRequest:   avgCostPerRequest,
			AvgLatency:          usage.AvgLatency,
			AvgTokensPerRequest: avgTokensPerRequest,
			EfficiencyScore:     efficiencyScore,
		}
	}

	return comparisons, nil
}

// Supporting types for analytics

// DashboardSummary contains all data for the main dashboard
type DashboardSummary struct {
	TimeRange     string                  `json:"time_range"`
	StartTime     time.Time               `json:"start_time"`
	EndTime       time.Time               `json:"end_time"`
	MetricSummary *models.MetricSummary   `json:"metric_summary"`
	CostBreakdown []*models.CostBreakdown `json:"cost_breakdown"`
	ModelUsage    []*models.ModelUsage    `json:"model_usage"`
	Insights      []Insight               `json:"insights"`
}

// Insight represents an actionable insight
type Insight struct {
	Type        string `json:"type"`     // success, info, warning, error
	Category    string `json:"category"` // cost, performance, reliability
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // info, low, medium, high, critical
}

// CostAnalysis contains detailed cost information
type CostAnalysis struct {
	TotalCost          float64                 `json:"total_cost"`
	DailyAverage       float64                 `json:"daily_average"`
	MonthlyProjection  float64                 `json:"monthly_projection"`
	CostBreakdown      []*models.CostBreakdown `json:"cost_breakdown"`
	MostExpensiveModel string                  `json:"most_expensive_model"`
	HighestCost        float64                 `json:"highest_cost"`
}

// PerformanceMetrics contains performance analysis
type PerformanceMetrics struct {
	Summary        *models.MetricSummary `json:"summary"`
	Status         string                `json:"status"` // excellent, good, fair, poor
	Recommendation string                `json:"recommendation"`
}

// ModelComparison compares different models
type ModelComparison struct {
	Model               string  `json:"model"`
	TotalCalls          int64   `json:"total_calls"`
	TotalCost           float64 `json:"total_cost"`
	AvgCostPerRequest   float64 `json:"avg_cost_per_request"`
	AvgLatency          float64 `json:"avg_latency"`
	AvgTokensPerRequest float64 `json:"avg_tokens_per_request"`
	EfficiencyScore     float64 `json:"efficiency_score"` // Lower is better
}
