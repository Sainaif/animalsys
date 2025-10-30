package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/pkg"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Extract token
		token, err := pkg.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		// Validate token
		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			if err == pkg.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
					"code": "TOKEN_EXPIRED",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// User info is set in context if token is present and valid, but doesn't abort if missing
func OptionalAuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := pkg.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// GetUserIDFromContext gets user ID from context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetUserRoleFromContext gets user role from context
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	roleStr, ok := role.(string)
	return roleStr, ok
}

// GetUserEmailFromContext gets user email from context
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}

// MustGetUserID gets user ID from context or panics
func MustGetUserID(c *gin.Context) string {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		panic("user_id not found in context")
	}
	return userID
}

// MustGetUserRole gets user role from context or panics
func MustGetUserRole(c *gin.Context) string {
	role, exists := GetUserRoleFromContext(c)
	if !exists {
		panic("user_role not found in context")
	}
	return role
}
