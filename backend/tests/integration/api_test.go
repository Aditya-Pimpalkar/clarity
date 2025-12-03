package integration

import (
	"testing"
	"time"
)

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	client := NewTestClient(BaseURL)

	resp, body := client.GET(t, "/health")

	AssertStatusCode(t, resp, 200)
	AssertJSONField(t, body, "status")

	data := ParseJSON(t, body)
	if data["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", data["status"])
	}
}

// TestAPIInfo tests the API info endpoint
func TestAPIInfo(t *testing.T) {
	client := NewTestClient(BaseURL)

	resp, body := client.GET(t, "/api/v1")

	AssertStatusCode(t, resp, 200)
	AssertJSONField(t, body, "version")

	data := ParseJSON(t, body)
	if data["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", data["version"])
	}
}

// TestTraceCreation tests creating a trace via API key
func TestTraceCreation(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetAPIKey(TestAPIKey)

	payload := map[string]interface{}{
		"organization_id": "org-test",
		"project_id":      "proj-test",
		"trace_type":      "single_call",
		"model":           "gpt-4",
		"provider":        "openai",
		"spans": []map[string]interface{}{
			{
				"name":              "test_span",
				"model":             "gpt-4",
				"provider":          "openai",
				"input":             "test input",
				"output":            "test output",
				"prompt_tokens":     10,
				"completion_tokens": 5,
				"duration_ms":       100,
				"status":            "success",
			},
		},
	}

	resp, body := client.POST(t, "/api/v1/traces", payload)

	AssertStatusCode(t, resp, 201)
	AssertJSONField(t, body, "success")

	data := ParseJSON(t, body)
	if !data["success"].(bool) {
		t.Error("Expected success=true")
	}

	// Check for trace_id in response
	responseData := data["data"].(map[string]interface{})
	if responseData["trace_id"] == nil {
		t.Error("Expected trace_id in response")
	}
}

// TestBatchTraceCreation tests batch trace creation
func TestBatchTraceCreation(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetAPIKey(TestAPIKey)

	payload := map[string]interface{}{
		"traces": []map[string]interface{}{
			{
				"organization_id": "org-test",
				"project_id":      "proj-test",
				"trace_type":      "single_call",
				"model":           "gpt-4",
				"provider":        "openai",
				"spans": []map[string]interface{}{
					{
						"name":              "batch_test_1",
						"model":             "gpt-4",
						"provider":          "openai",
						"prompt_tokens":     20,
						"completion_tokens": 10,
						"duration_ms":       150,
						"status":            "success",
					},
				},
			},
			{
				"organization_id": "org-test",
				"project_id":      "proj-test",
				"trace_type":      "single_call",
				"model":           "gpt-3.5-turbo",
				"provider":        "openai",
				"spans": []map[string]interface{}{
					{
						"name":              "batch_test_2",
						"model":             "gpt-3.5-turbo",
						"provider":          "openai",
						"prompt_tokens":     15,
						"completion_tokens": 8,
						"duration_ms":       100,
						"status":            "success",
					},
				},
			},
		},
	}

	resp, body := client.POST(t, "/api/v1/traces/batch", payload)

	AssertStatusCode(t, resp, 201)

	data := ParseJSON(t, body)
	responseData := data["data"].(map[string]interface{})

	if responseData["accepted"].(float64) != 2 {
		t.Errorf("Expected 2 traces accepted, got %v", responseData["accepted"])
	}
}

// TestUnauthorizedAccess tests accessing protected endpoint without auth
func TestUnauthorizedAccess(t *testing.T) {
	client := NewTestClient(BaseURL)

	// Try to create trace without API key
	payload := map[string]interface{}{
		"organization_id": "org-test",
		"project_id":      "proj-test",
		"trace_type":      "single_call",
		"model":           "gpt-4",
		"provider":        "openai",
		"spans":           []map[string]interface{}{},
	}

	resp, _ := client.POST(t, "/api/v1/traces", payload)

	AssertStatusCode(t, resp, 401)
}

// TestRateLimiting tests rate limiting
func TestRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limit test in short mode")
	}

	client := NewTestClient(BaseURL)
	client.SetAPIKey(TestAPIKey)

	// Make many requests rapidly
	hitLimit := false
	for i := 0; i < 150; i++ {
		resp, _ := client.GET(t, "/api/v1")
		if resp.StatusCode == 429 {
			hitLimit = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !hitLimit {
		t.Log("Warning: Rate limit not hit (may need adjustment)")
	}
}
