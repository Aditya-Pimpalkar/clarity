package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sahared/llm-observability/internal/services"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboard handles GET /api/v1/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	// Parse query parameters
	orgID := c.Query("organization_id")
	projectID := c.Query("project_id")
	timeRange := c.Query("time_range", "24h")

	// Validate required parameters
	if orgID == "" {
		return BadRequestResponse(c, "organization_id is required")
	}

	if projectID == "" {
		return BadRequestResponse(c, "project_id is required")
	}

	// Call service
	summary, err := h.analyticsService.GetDashboardSummary(c.Context(), orgID, projectID, timeRange)
	if err != nil {
		return InternalErrorResponse(c, "Failed to get dashboard summary: "+err.Error())
	}

	// Return response
	return SuccessResponse(c, summary)
}

// GetCostAnalysis handles GET /api/v1/analytics/costs
func (h *AnalyticsHandler) GetCostAnalysis(c *fiber.Ctx) error {
	// Parse query parameters
	orgID := c.Query("organization_id")
	projectID := c.Query("project_id")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// Validate required parameters
	if orgID == "" {
		return BadRequestResponse(c, "organization_id is required")
	}

	if projectID == "" {
		return BadRequestResponse(c, "project_id is required")
	}

	// Parse time range
	var startTime, endTime time.Time
	var err error

	if startTimeStr == "" {
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	} else {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return BadRequestResponse(c, "Invalid start_time format")
		}
	}

	if endTimeStr == "" {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return BadRequestResponse(c, "Invalid end_time format")
		}
	}

	// Call service
	analysis, err := h.analyticsService.GetCostAnalysis(c.Context(), orgID, projectID, startTime, endTime)
	if err != nil {
		return InternalErrorResponse(c, "Failed to get cost analysis: "+err.Error())
	}

	// Return response
	return SuccessResponse(c, analysis)
}

// GetPerformanceMetrics handles GET /api/v1/analytics/performance
func (h *AnalyticsHandler) GetPerformanceMetrics(c *fiber.Ctx) error {
	orgID := c.Query("organization_id")
	projectID := c.Query("project_id")

	if orgID == "" || projectID == "" {
		return BadRequestResponse(c, "organization_id and project_id are required")
	}

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	metrics, err := h.analyticsService.GetPerformanceMetrics(c.Context(), orgID, projectID, startTime, endTime)
	if err != nil {
		return InternalErrorResponse(c, "Failed to get performance metrics: "+err.Error())
	}

	return SuccessResponse(c, metrics)
}

// GetModelComparison handles GET /api/v1/analytics/models
func (h *AnalyticsHandler) GetModelComparison(c *fiber.Ctx) error {
	orgID := c.Query("organization_id")
	projectID := c.Query("project_id")

	if orgID == "" || projectID == "" {
		return BadRequestResponse(c, "organization_id and project_id are required")
	}

	startTime := time.Now().Add(-7 * 24 * time.Hour)
	endTime := time.Now()

	comparison, err := h.analyticsService.GetModelComparison(c.Context(), orgID, projectID, startTime, endTime)
	if err != nil {
		return InternalErrorResponse(c, "Failed to get model comparison: "+err.Error())
	}

	return SuccessResponse(c, comparison)
}
