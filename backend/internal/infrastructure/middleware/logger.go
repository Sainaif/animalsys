package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Get error if any
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Build query string
		if raw != "" {
			path = path + "?" + raw
		}

		// Log request
		logEvent := log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Str("ip", clientIP).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent())

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			logEvent = logEvent.Str("user_id", userID.(string))
		}

		// Add error if exists
		if errorMessage != "" {
			logEvent = logEvent.Str("error", errorMessage)
		}

		// Determine log level based on status code
		if statusCode >= 500 {
			logEvent = log.Error().
				Str("method", c.Request.Method).
				Str("path", path).
				Int("status", statusCode).
				Str("ip", clientIP).
				Dur("latency", latency).
				Str("error", errorMessage)
		} else if statusCode >= 400 {
			logEvent = log.Warn().
				Str("method", c.Request.Method).
				Str("path", path).
				Int("status", statusCode).
				Str("ip", clientIP).
				Dur("latency", latency)
		}

		logEvent.Msg("HTTP Request")
	}
}

// LoggerWithSkipPaths logs HTTP requests but skips certain paths
func LoggerWithSkipPaths(skipPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip logging for specified paths
		for _, skipPath := range skipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// Use regular logger
		Logger()(c)
	}
}
