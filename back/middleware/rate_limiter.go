package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// RateLimiter represents a rate limiter for a specific IP
type RateLimiter struct {
	requests []time.Time
	mutex    sync.Mutex
}

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	MaxRequests int           // Maximum requests allowed
	Window      time.Duration // Time window for rate limiting
	Message     string        // Custom message for rate limit exceeded
}

// DefaultRateLimiterConfig provides default configuration
var DefaultRateLimiterConfig = RateLimiterConfig{
	MaxRequests: 100,
	Window:      time.Minute,
	Message:     "Rate limit exceeded. Please try again later.",
}

// Global map to store rate limiters for each IP
var rateLimiters = make(map[string]*RateLimiter)
var rateLimitersMutex sync.RWMutex

// RateLimit creates a rate limiting middleware
func RateLimit(config ...RateLimiterConfig) fiber.Handler {
	cfg := DefaultRateLimiterConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c fiber.Ctx) error {
		ip := c.IP()

		// Get or create rate limiter for this IP
		rateLimitersMutex.RLock()
		limiter, exists := rateLimiters[ip]
		rateLimitersMutex.RUnlock()

		if !exists {
			limiter = &RateLimiter{
				requests: make([]time.Time, 0),
			}
			rateLimitersMutex.Lock()
			rateLimiters[ip] = limiter
			rateLimitersMutex.Unlock()
		}

		// Check rate limit
		if !limiter.Allow(cfg.MaxRequests, cfg.Window) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       cfg.Message,
				"retry_after": cfg.Window.Seconds(),
			})
		}

		return c.Next()
	}
}

// Allow checks if a request is allowed based on rate limiting rules
func (rl *RateLimiter) Allow(maxRequests int, window time.Duration) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Remove old requests outside the window
	validRequests := make([]time.Time, 0)
	for _, reqTime := range rl.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	rl.requests = validRequests

	// Check if we can allow this request
	if len(rl.requests) >= maxRequests {
		return false
	}

	// Add current request
	rl.requests = append(rl.requests, now)
	return true
}

// StrictRateLimit creates a stricter rate limiting middleware for sensitive endpoints
func StrictRateLimit() fiber.Handler {
	return RateLimit(RateLimiterConfig{
		MaxRequests: 10,
		Window:      time.Minute,
		Message:     "Too many requests to sensitive endpoint. Please try again later.",
	})
}

// AuthRateLimit creates rate limiting specifically for authentication endpoints
func AuthRateLimit() fiber.Handler {
	return RateLimit(RateLimiterConfig{
		MaxRequests: 5,
		Window:      time.Minute * 5,
		Message:     "Too many authentication attempts. Please try again in 5 minutes.",
	})
}

// CleanupRateLimiters periodically cleans up old rate limiters to prevent memory leaks
func CleanupRateLimiters() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			rateLimitersMutex.Lock()
			now := time.Now()
			for ip, limiter := range rateLimiters {
				limiter.mutex.Lock()
				// Remove rate limiter if no requests in the last hour
				if len(limiter.requests) == 0 || now.Sub(limiter.requests[len(limiter.requests)-1]) > time.Hour {
					delete(rateLimiters, ip)
				}
				limiter.mutex.Unlock()
			}
			rateLimitersMutex.Unlock()
		}
	}()
}
