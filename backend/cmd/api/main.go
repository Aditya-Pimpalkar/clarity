package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/sahared/llm-observability/internal/api"
	"github.com/sahared/llm-observability/internal/repository"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	}

	// Build ClickHouse DSN
	dsn := buildClickHouseDSN()

	// Connect to ClickHouse
	log.Println("🔌 Connecting to ClickHouse...")
	repo, err := repository.NewClickHouseRepository(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to ClickHouse:", err)
	}
	defer repo.Close()
	log.Println("✅ Connected to ClickHouse")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "LLM Observability Platform API",
		ServerHeader: "LLM-Observability",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} (${latency}) [${locals:requestid}]\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "LLM Observability Platform",
			"version": "1.0.0",
			"status":  "running",
			"message": "Welcome to LLM Observability Platform API",
		})
	})

	// Setup all API routes (this includes /health, /ready, /live, and /api/v1/*)
	api.SetupRoutes(app, repo)

	// Start server
	port := getEnv("PORT", "8080")

	log.Println(strings.Repeat("=", 60))
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 Environment: %s", getEnv("ENV", "development"))
	log.Printf("🌐 Root: http://localhost:%s/", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)
	log.Printf("📊 API v1: http://localhost:%s/api/v1", port)
	log.Printf("📈 Dashboard: http://localhost:%s/api/v1/analytics/dashboard", port)
	log.Println(strings.Repeat("=", 60))

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// buildClickHouseDSN builds the ClickHouse connection string from environment variables
func buildClickHouseDSN() string {
	host := getEnv("CLICKHOUSE_HOST", "localhost")
	port := getEnv("CLICKHOUSE_PORT", "9000")
	database := getEnv("CLICKHOUSE_DATABASE", "llm_observability")
	user := getEnv("CLICKHOUSE_USER", "default")
	password := getEnv("CLICKHOUSE_PASSWORD", "")

	dsn := "clickhouse://"
	if user != "" {
		dsn += user
		if password != "" {
			dsn += ":" + password
		}
		dsn += "@"
	}
	dsn += host + ":" + port + "/" + database

	log.Printf("📡 ClickHouse DSN: clickhouse://%s:%s/%s", host, port, database)
	return dsn
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
