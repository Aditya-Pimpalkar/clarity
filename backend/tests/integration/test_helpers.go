package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestClient is a helper for making API requests in tests
type TestClient struct {
	BaseURL string
	APIKey  string
	Token   string
	Client  *http.Client
}

// NewTestClient creates a new test client
func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetAPIKey sets the API key for requests
func (c *TestClient) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

// SetToken sets the JWT token for requests
func (c *TestClient) SetToken(token string) {
	c.Token = token
}

// GET makes a GET request
func (c *TestClient) GET(t *testing.T, path string) (*http.Response, []byte) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	c.setHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return resp, body
}

// POST makes a POST request with JSON body
func (c *TestClient) POST(t *testing.T, path string, payload interface{}) (*http.Response, []byte) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.setHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return resp, body
}

// setHeaders sets authentication headers
func (c *TestClient) setHeaders(req *http.Request) {
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
}

// AssertStatusCode asserts the response status code
func AssertStatusCode(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Errorf("Expected status %d, got %d", expected, resp.StatusCode)
	}
}

// AssertJSONField asserts a field exists in JSON response
func AssertJSONField(t *testing.T, body []byte, field string) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, exists := data[field]; !exists {
		t.Errorf("Expected field '%s' not found in response", field)
	}
}

// ParseJSON parses JSON response
func ParseJSON(t *testing.T, body []byte) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	return data
}

// WaitForServer waits for server to be ready
func WaitForServer(baseURL string, maxAttempts int) error {
	client := &http.Client{Timeout: 2 * time.Second}

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("server not ready after %d attempts", maxAttempts)
}
