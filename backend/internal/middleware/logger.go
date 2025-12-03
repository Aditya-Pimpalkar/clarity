package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID)

		start := time.Now()

		fmt.Printf("[%s] --> %s %s | Request ID: %s\n",
			start.Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			requestID,
		)

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		statusColor := getStatusColor(statusCode)

		fmt.Printf("[%s] <-- %s %s | %s%d%s | %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			statusColor,
			statusCode,
			"\033[0m",
			duration,
		)

		return err
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m"
	case status >= 300 && status < 400:
		return "\033[36m"
	case status >= 400 && status < 500:
		return "\033[33m"
	case status >= 500:
		return "\033[31m"
	default:
		return "\033[0m"
	}
}

func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Locals("request_id").(string)
	if !ok {
		return ""
	}
	return requestID
}
