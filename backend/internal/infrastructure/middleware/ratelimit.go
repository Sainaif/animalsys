package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

// Visitor represents a client with rate limit tracking
type Visitor struct {
	limiter  *TokenBucket
	lastSeen time.Time
}

// TokenBucket implements token bucket algorithm
type TokenBucket struct {
	tokens     int
	capacity   int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerWindow int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     requestsPerWindow,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// RateLimit returns a rate limiting middleware
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Check rate limit
		if !rl.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

// allow checks if the request is allowed
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[key]
	if !exists {
		// Create new visitor
		visitor = &Visitor{
			limiter: &TokenBucket{
				tokens:     rl.rate,
				capacity:   rl.rate,
				refillRate: rl.rate,
				lastRefill: time.Now(),
			},
			lastSeen: time.Now(),
		}
		rl.visitors[key] = visitor
	}

	// Update last seen
	visitor.lastSeen = time.Now()

	// Try to consume a token
	return visitor.limiter.consume()
}

// consume tries to consume one token
func (tb *TokenBucket) consume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Minutes()) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// Try to consume a token
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// cleanupVisitors removes old visitors periodically
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for key, visitor := range rl.visitors {
			// Remove visitors not seen in the last hour
			if time.Since(visitor.lastSeen) > time.Hour {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IPRateLimit returns a simple IP-based rate limiter
func IPRateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)
	return limiter.RateLimit()
}

// UserRateLimit returns a user-based rate limiter
func UserRateLimit(requestsPerHour int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerHour, time.Hour)

	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// If no user ID, fall back to IP-based limiting
			ip := c.ClientIP()
			if !limiter.allow(ip) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
					"message": "Too many requests. Please try again later.",
				})
				return
			}
			c.Next()
			return
		}

		// Check user-based rate limit
		if !limiter.allow(userID.(string)) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
