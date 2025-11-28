package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "LLM Observability Platform API",
		ServerHeader: "LLM-Observability",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "LLM Observability Platform",
			"version": "1.0.0",
			"status":  "running",
			"message": "Welcome to LLM Observability Platform API",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "llm-observability",
			"version": "1.0.0",
			"uptime":  time.Now().Unix(),
		})
	})

	v1 := app.Group("/api/v1")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "LLM Observability Platform API v1",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /        - Root endpoint",
				"GET /health  - Health check",
				"GET /api/v1  - API information",
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 Environment: %s", os.Getenv("ENV"))
	log.Printf("🌐 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API docs: http://localhost:%s/api/v1", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
