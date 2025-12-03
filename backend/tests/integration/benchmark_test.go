package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// BenchmarkHealthEndpoint benchmarks the health endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "http://localhost:8080/health"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(url)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 429 {
			b.Errorf("Unexpected status: %d", resp.StatusCode)
		}
	}
}

// BenchmarkTraceCreation benchmarks trace creation
func BenchmarkTraceCreation(b *testing.B) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "http://localhost:8080/api/v1/traces"

	payload := map[string]interface{}{
		"organization_id": "org-bench",
		"project_id":      "proj-bench",
		"trace_type":      "single_call",
		"model":           "gpt-4",
		"provider":        "openai",
		"spans": []map[string]interface{}{
			{
				"name":              "benchmark",
				"model":             "gpt-4",
				"provider":          "openai",
				"prompt_tokens":     10,
				"completion_tokens": 5,
				"duration_ms":       100,
				"status":            "success",
			},
		},
	}

	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", TestAPIKey)

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != 201 && resp.StatusCode != 429 {
			b.Errorf("Unexpected status: %d", resp.StatusCode)
		}
	}
}
