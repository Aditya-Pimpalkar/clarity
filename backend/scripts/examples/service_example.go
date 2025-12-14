package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

func main() {
	fmt.Println("🚀 LLM Observability Service Layer Example")
	fmt.Println("==========================================")

	// Step 1: Connect to ClickHouse
	fmt.Println("\n1. Connecting to ClickHouse...")
	dsn := "clickhouse://localhost:9000/llm_observability"
	repo, err := repository.NewClickHouseRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer repo.Close()
	fmt.Println("✓ Connected successfully")

	// Step 2: Create services
	fmt.Println("\n2. Creating services...")
	traceService := services.NewTraceService(repo)
	analyticsService := services.NewAnalyticsService(repo)
	userService := services.NewUserService(repo)
	fmt.Println("✓ Services created")

	ctx := context.Background()

	// Step 3: Create organization and project
	fmt.Println("\n3. Setting up organization and project...")
	org, err := userService.CreateOrganization(ctx, "Example Corp", "pro")
	if err != nil {
		log.Printf("Organization may already exist: %v", err)
		// Use a test org ID if creation fails
		org = &models.Organization{ID: "org-example"}
	} else {
		fmt.Printf("✓ Created organization: %s\n", org.ID)
	}

	project, err := userService.CreateProject(ctx, org.ID, "Example Project", "Demo project for testing")
	if err != nil {
		log.Printf("Project may already exist: %v", err)
		project = &models.Project{ID: "proj-example", OrganizationID: org.ID}
	} else {
		fmt.Printf("✓ Created project: %s\n", project.ID)
	}

	// Step 4: Create a trace using the service
	fmt.Println("\n4. Creating a trace...")
	traceReq := &models.TraceRequest{
		OrganizationID: org.ID,
		ProjectID:      project.ID,
		TraceType:      "single_call",
		Model:          "gpt-4",
		Provider:       "openai",
		UserID:         "user-example",
		Metadata: map[string]interface{}{
			"environment": "development",
			"version":     "1.0.0",
		},
		Spans: []models.SpanRequest{
			{
				Name:             "chat_completion",
				Model:            "gpt-4",
				Provider:         "openai",
				Input:            "What is the capital of France?",
				Output:           "The capital of France is Paris.",
				PromptTokens:     15,
				CompletionTokens: 8,
				DurationMs:       245,
				Status:           "success",
				Metadata: map[string]interface{}{
					"temperature": 0.7,
				},
			},
		},
	}

	resp, err := traceService.CreateTrace(ctx, traceReq)
	if err != nil {
		log.Fatalf("Failed to create trace: %v", err)
	}
	fmt.Printf("✓ Trace created: %s\n", resp.TraceID)
	fmt.Printf("  Status: %s\n", resp.Status)
	fmt.Printf("  Created at: %s\n", resp.CreatedAt)

	// Small delay to let metrics be recorded
	time.Sleep(100 * time.Millisecond)

	// Step 5: Create more sample traces
	fmt.Println("\n5. Creating additional sample traces...")
	for i := 0; i < 4; i++ {
		traceReq := &models.TraceRequest{
			OrganizationID: org.ID,
			ProjectID:      project.ID,
			TraceType:      "single_call",
			Model:          []string{"gpt-4", "gpt-3.5-turbo", "claude-3-sonnet"}[i%3],
			Provider:       []string{"openai", "openai", "anthropic"}[i%3],
			UserID:         "user-example",
			Spans: []models.SpanRequest{
				{
					Name:             "llm_call",
					Model:            []string{"gpt-4", "gpt-3.5-turbo", "claude-3-sonnet"}[i%3],
					Provider:         []string{"openai", "openai", "anthropic"}[i%3],
					Input:            "Sample input",
					Output:           "Sample output",
					PromptTokens:     uint32(50 + i*10),
					CompletionTokens: uint32(30 + i*5),
					DurationMs:       uint32(150 + i*50),
					Status:           "success",
				},
			},
		}
		_, err := traceService.CreateTrace(ctx, traceReq)
		if err != nil {
			log.Printf("Warning: Failed to create trace %d: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("✓ Created 4 additional traces")

	// Step 6: Get dashboard summary
	fmt.Println("\n6. Getting dashboard summary...")
	summary, err := analyticsService.GetDashboardSummary(ctx, org.ID, project.ID, "24h")
	if err != nil {
		log.Fatalf("Failed to get dashboard summary: %v", err)
	}

	fmt.Println("✓ Dashboard Summary:")
	fmt.Printf("  Total Requests: %d\n", summary.MetricSummary.TotalRequests)
	fmt.Printf("  Total Cost: $%.4f\n", summary.MetricSummary.TotalCostUSD)
	fmt.Printf("  Avg Latency: %.2fms\n", summary.MetricSummary.AvgLatencyMs)
	fmt.Printf("  Success Rate: %.2f%%\n", summary.MetricSummary.SuccessRate)

	if len(summary.CostBreakdown) > 0 {
		fmt.Println("\n  Cost Breakdown:")
		for i, cb := range summary.CostBreakdown {
			fmt.Printf("    %d. %s: $%.4f (%.1f%%)\n", i+1, cb.Model, cb.TotalCost, cb.Percentage)
		}
	}

	if len(summary.Insights) > 0 {
		fmt.Println("\n  Insights:")
		for i, insight := range summary.Insights {
			fmt.Printf("    %d. [%s] %s\n", i+1, insight.Type, insight.Title)
			fmt.Printf("       %s\n", insight.Description)
		}
	}

	// Step 7: Get model comparison
	fmt.Println("\n7. Getting model comparison...")
	comparison, err := analyticsService.GetModelComparison(ctx, org.ID, project.ID,
		time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		log.Fatalf("Failed to get model comparison: %v", err)
	}

	if len(comparison) > 0 {
		fmt.Println("✓ Model Comparison:")
		for i, comp := range comparison {
			fmt.Printf("  %d. %s\n", i+1, comp.Model)
			fmt.Printf("     Calls: %d\n", comp.TotalCalls)
			fmt.Printf("     Avg Cost: $%.4f\n", comp.AvgCostPerRequest)
			fmt.Printf("     Avg Latency: %.2fms\n", comp.AvgLatency)
			fmt.Printf("     Efficiency Score: %.2f\n", comp.EfficiencyScore)
		}
	}

	fmt.Println("\n✓ All examples completed successfully!")
}
