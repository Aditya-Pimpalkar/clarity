package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSConfig() fiber.Handler {
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:8080"
	}

	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		}, ","),
		AllowCredentials: true,
		ExposeHeaders: strings.Join([]string{
			"X-Request-ID",
			"X-Total-Count",
		}, ","),
		MaxAge: 3600,
	})
}
