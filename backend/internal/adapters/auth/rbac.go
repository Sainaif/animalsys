package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role represents user roles
type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleEmployee   Role = "employee"
	RoleVolunteer  Role = "volunteer"
	RoleUser       Role = "user"
	RoleGuest      Role = "guest"
)

// RoleHierarchy defines the hierarchy of roles (higher roles include permissions of lower roles)
var RoleHierarchy = map[Role]int{
	RoleSuperAdmin: 6,
	RoleAdmin:      5,
	RoleEmployee:   4,
	RoleVolunteer:  3,
	RoleUser:       2,
	RoleGuest:      1,
}

// RequireRole creates a middleware that requires specific role(s)
func RequireRole(allowedRoles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetUserRoleFromContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Check if user has any of the allowed roles
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if Role(userRole) == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"required_roles": allowedRoles,
			})
			return
		}

		c.Next()
	}
}

// RequireMinRole creates a middleware that requires minimum role level
func RequireMinRole(minRole Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetUserRoleFromContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		userRoleLevel, ok := RoleHierarchy[Role(userRole)]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid role",
			})
			return
		}

		minRoleLevel := RoleHierarchy[minRole]
		if userRoleLevel < minRoleLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"required_min_role": minRole,
			})
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin requires super admin role
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(RoleSuperAdmin)
}

// RequireAdmin requires admin or super admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireMinRole(RoleAdmin)
}

// RequireEmployee requires employee, admin, or super admin role
func RequireEmployee() gin.HandlerFunc {
	return RequireMinRole(RoleEmployee)
}

// RequireVolunteer requires volunteer or higher role
func RequireVolunteer() gin.HandlerFunc {
	return RequireMinRole(RoleVolunteer)
}

// RequireAuthenticated requires any authenticated user
func RequireAuthenticated() gin.HandlerFunc {
	return RequireMinRole(RoleUser)
}

// CheckPermission checks if user has permission for a specific action
func CheckPermission(userRole Role, requiredRole Role) bool {
	userLevel, userExists := RoleHierarchy[userRole]
	requiredLevel, requiredExists := RoleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// CanAccessResource checks if user can access a specific resource
func CanAccessResource(c *gin.Context, resourceOwnerID string) bool {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		return false
	}

	userRole, exists := GetUserRoleFromContext(c)
	if !exists {
		return false
	}

	// Super admin and admin can access all resources
	if Role(userRole) == RoleSuperAdmin || Role(userRole) == RoleAdmin {
		return true
	}

	// User can access their own resources
	if userID == resourceOwnerID {
		return true
	}

	return false
}

// RequireResourceOwnership requires user to be the owner or admin
func RequireResourceOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		userRole, exists := GetUserRoleFromContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Admin and super admin can access everything
		if Role(userRole) == RoleSuperAdmin || Role(userRole) == RoleAdmin {
			c.Next()
			return
		}

		// For other users, check resource ownership
		// Resource owner ID should be set in the route handler
		resourceOwnerID, exists := c.Get("resource_owner_id")
		if !exists {
			// If not set, assume user is trying to access their own resource
			c.Next()
			return
		}

		if userID != resourceOwnerID.(string) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have permission to access this resource",
			})
			return
		}

		c.Next()
	}
}

// IsRole checks if user has specific role
func IsRole(c *gin.Context, role Role) bool {
	userRole, exists := GetUserRoleFromContext(c)
	if !exists {
		return false
	}

	return Role(userRole) == role
}

// HasMinRole checks if user has minimum role level
func HasMinRole(c *gin.Context, minRole Role) bool {
	userRole, exists := GetUserRoleFromContext(c)
	if !exists {
		return false
	}

	userLevel, ok := RoleHierarchy[Role(userRole)]
	if !ok {
		return false
	}

	minLevel := RoleHierarchy[minRole]
	return userLevel >= minLevel
}
