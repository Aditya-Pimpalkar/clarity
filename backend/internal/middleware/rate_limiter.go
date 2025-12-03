package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]*clientInfo
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type clientInfo struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientInfo),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier (IP address or user ID)
		clientID := c.IP()

		// If authenticated, use user ID for more accurate limiting
		if userID := GetUserID(c); userID != "" {
			clientID = userID
		}

		// Check rate limit
		allowed, resetTime := rl.allow(clientID)
		if !allowed {
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": time.Until(resetTime).Seconds(),
			})
		}

		// Set rate limit headers
		remaining := rl.getRemaining(clientID)
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		return c.Next()
	}
}

// allow checks if a request is allowed
func (rl *RateLimiter) allow(clientID string) (bool, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	info, exists := rl.requests[clientID]
	if !exists {
		rl.requests[clientID] = &clientInfo{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true, now.Add(rl.window)
	}

	// Check if window has expired
	if now.After(info.resetTime) {
		info.count = 1
		info.resetTime = now.Add(rl.window)
		return true, info.resetTime
	}

	// Check if limit exceeded
	if info.count >= rl.limit {
		return false, info.resetTime
	}

	info.count++
	return true, info.resetTime
}

// getRemaining returns the number of remaining requests
func (rl *RateLimiter) getRemaining(clientID string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	info, exists := rl.requests[clientID]
	if !exists {
		return rl.limit
	}

	remaining := rl.limit - info.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanup periodically removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientID, info := range rl.requests {
			if now.After(info.resetTime.Add(rl.window)) {
				delete(rl.requests, clientID)
			}
		}
		rl.mu.Unlock()
	}
}

// GlobalRateLimiter creates a rate limiter for all requests
// 1000 requests per hour
func GlobalRateLimiter() fiber.Handler {
	limiter := NewRateLimiter(1000, time.Hour)
	return limiter.Middleware()
}

// StrictRateLimiter creates a strict rate limiter
// 100 requests per minute
func StrictRateLimiter() fiber.Handler {
	limiter := NewRateLimiter(100, time.Minute)
	return limiter.Middleware()
}

// APIKeyRateLimiter creates a rate limiter for API keys
// 10000 requests per hour
func APIKeyRateLimiter() fiber.Handler {
	limiter := NewRateLimiter(10000, time.Hour)
	return limiter.Middleware()
}
