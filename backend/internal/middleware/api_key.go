package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth middleware for API key authentication
// Used by SDK clients for programmatic access
func APIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get API key from header
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header for API keys
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing API key",
				"code":  "MISSING_API_KEY",
				"hint":  "Provide X-API-Key header or Authorization: ApiKey <key>",
			})
		}

		// Validate API key
		orgID, projectID, valid := validateAPIKey(apiKey)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
		}

		// Store in context for later use
		c.Locals(string(OrgIDKey), orgID)
		c.Locals(string(ProjectIDKey), projectID)
		c.Locals("authenticated_by", "api_key")

		return c.Next()
	}
}

// validateAPIKey validates an API key and returns org/project IDs
func validateAPIKey(apiKey string) (orgID, projectID string, valid bool) {
	// Get test keys
	testKeys := getTestAPIKeys()

	// Check direct key match first (for development)
	if info, exists := testKeys[apiKey]; exists {
		return info.OrgID, info.ProjectID, true
	}

	// Also check hashed version
	hash := sha256.Sum256([]byte(apiKey))
	keyHash := hex.EncodeToString(hash[:])

	if info, exists := testKeys[keyHash]; exists {
		return info.OrgID, info.ProjectID, true
	}

	return "", "", false
}

// APIKeyInfo stores API key metadata
type APIKeyInfo struct {
	OrgID     string
	ProjectID string
}

// getTestAPIKeys returns test API keys for development
func getTestAPIKeys() map[string]APIKeyInfo {
	return map[string]APIKeyInfo{
		// Direct key (for easy testing)
		"test-key-123": {
			OrgID:     "org-test",
			ProjectID: "proj-test",
		},
		"demo-key-456": {
			OrgID:     "org-demo",
			ProjectID: "proj-demo",
		},
		// Hashed version of "test-key-123"
		"ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae": {
			OrgID:     "org-test",
			ProjectID: "proj-test",
		},
	}
}

// OptionalAuth allows requests with or without authentication
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try JWT first
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			authMiddleware := AuthMiddleware()
			if err := authMiddleware(c); err != nil {
				return c.Next()
			}
			return c.Next()
		}

		// Try API key
		apiKey := c.Get("X-API-Key")
		if apiKey != "" || strings.HasPrefix(authHeader, "ApiKey ") {
			apiKeyMiddleware := APIKeyAuth()
			if err := apiKeyMiddleware(c); err != nil {
				return c.Next()
			}
			return c.Next()
		}

		return c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := GetRole(c)
		if role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role information not found",
				"code":  "NO_ROLE",
			})
		}

		// Role hierarchy: admin > member > viewer
		roleLevel := map[string]int{
			"viewer": 1,
			"member": 2,
			"admin":  3,
		}

		userLevel := roleLevel[role]
		requiredLevel := roleLevel[requiredRole]

		if userLevel < requiredLevel {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":         "Insufficient permissions",
				"code":          "INSUFFICIENT_PERMISSIONS",
				"required_role": requiredRole,
				"user_role":     role,
			})
		}

		return c.Next()
	}
}
