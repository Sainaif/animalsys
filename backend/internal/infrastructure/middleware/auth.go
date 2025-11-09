package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/sainaif/animalsys/backend/pkg/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware validates JWT tokens and attaches user to context
func AuthMiddleware(jwtService *security.JWTService, userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Get user ID from claims
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			c.Abort()
			return
		}

		// Fetch user from database
		user, err := userRepo.FindByID(c.Request.Context(), userID)
		if err != nil {
			if err == errors.ErrNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			}
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive() {
			c.JSON(http.StatusForbidden, gin.H{"error": "user account is not active"})
			c.Abort()
			return
		}

		// Attach user to context
		c.Set("user", user)
		c.Set("user_id", userID)

		c.Next()
	}
}

// OptionalAuth middleware that doesn't require authentication but attaches user if token is provided
func OptionalAuth(jwtService *security.JWTService, userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		user, err := userRepo.FindByID(c.Request.Context(), userID)
		if err != nil {
			c.Next()
			return
		}

		if user.IsActive() {
			c.Set("user", user)
			c.Set("user_id", userID)
		}

		c.Next()
	}
}

// GetUserFromContext retrieves the user from the Gin context
func GetUserFromContext(c *gin.Context) (*primitive.ObjectID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, errors.ErrUnauthorized
	}

	id, ok := userID.(primitive.ObjectID)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	return &id, nil
}
