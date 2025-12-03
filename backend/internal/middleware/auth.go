package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	OrgIDKey     contextKey = "org_id"
	RoleKey      contextKey = "role"
	ProjectIDKey contextKey = "project_id"
)

type Claims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
				"code":  "UNAUTHORIZED",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
				"code":  "INVALID_AUTH_FORMAT",
			})
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-secret-key-change-in-production"
			}
			return []byte(secret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
				"code":  "INVALID_CLAIMS",
			})
		}

		c.Locals(string(UserIDKey), claims.UserID)
		c.Locals(string(OrgIDKey), claims.OrganizationID)
		c.Locals(string(RoleKey), claims.Role)

		return c.Next()
	}
}

func GenerateToken(userID, organizationID, role, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:         userID,
		OrganizationID: organizationID,
		Role:           role,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "llm-observability",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func GetUserID(c *fiber.Ctx) string {
	userID, _ := c.Locals(string(UserIDKey)).(string)
	return userID
}

func GetOrgID(c *fiber.Ctx) string {
	orgID, _ := c.Locals(string(OrgIDKey)).(string)
	return orgID
}

func GetRole(c *fiber.Ctx) string {
	role, _ := c.Locals(string(RoleKey)).(string)
	return role
}
