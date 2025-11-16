package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/pkg/errors"
)

// Permission represents a specific permission
type Permission string

const (
	// User permissions
	PermissionViewUsers   Permission = "users:view"
	PermissionCreateUsers Permission = "users:create"
	PermissionUpdateUsers Permission = "users:update"
	PermissionDeleteUsers Permission = "users:delete"

	// Animal permissions
	PermissionViewAnimals   Permission = "animals:view"
	PermissionCreateAnimals Permission = "animals:create"
	PermissionUpdateAnimals Permission = "animals:update"
	PermissionDeleteAnimals Permission = "animals:delete"

	// Adoption permissions
	PermissionViewAdoptions   Permission = "adoptions:view"
	PermissionCreateAdoptions Permission = "adoptions:create"
	PermissionUpdateAdoptions Permission = "adoptions:update"
	PermissionDeleteAdoptions Permission = "adoptions:delete"

	// Veterinary permissions
	PermissionViewVeterinary   Permission = "veterinary:view"
	PermissionCreateVeterinary Permission = "veterinary:create"
	PermissionUpdateVeterinary Permission = "veterinary:update"
	PermissionDeleteVeterinary Permission = "veterinary:delete"

	// Donor permissions
	PermissionViewDonors   Permission = "donors:view"
	PermissionCreateDonors Permission = "donors:create"
	PermissionUpdateDonors Permission = "donors:update"
	PermissionDeleteDonors Permission = "donors:delete"

	// Donation permissions
	PermissionViewDonations   Permission = "donations:view"
	PermissionCreateDonations Permission = "donations:create"
	PermissionUpdateDonations Permission = "donations:update"
	PermissionDeleteDonations Permission = "donations:delete"

	// Campaign permissions
	PermissionViewCampaigns   Permission = "campaigns:view"
	PermissionCreateCampaigns Permission = "campaigns:create"
	PermissionUpdateCampaigns Permission = "campaigns:update"
	PermissionDeleteCampaigns Permission = "campaigns:delete"

	// Event permissions
	PermissionViewEvents   Permission = "events:view"
	PermissionCreateEvents Permission = "events:create"
	PermissionUpdateEvents Permission = "events:update"
	PermissionDeleteEvents Permission = "events:delete"

	// Volunteer permissions
	PermissionViewVolunteers   Permission = "volunteers:view"
	PermissionCreateVolunteers Permission = "volunteers:create"
	PermissionUpdateVolunteers Permission = "volunteers:update"
	PermissionDeleteVolunteers Permission = "volunteers:delete"

	// Communication permissions
	PermissionViewCommunications   Permission = "communications:view"
	PermissionCreateCommunications Permission = "communications:create"
	PermissionUpdateCommunications Permission = "communications:update"
	PermissionDeleteCommunications Permission = "communications:delete"

	// Template permissions
	PermissionViewTemplates   Permission = "templates:view"
	PermissionCreateTemplates Permission = "templates:create"
	PermissionUpdateTemplates Permission = "templates:update"
	PermissionDeleteTemplates Permission = "templates:delete"

	// Notification permissions
	PermissionViewNotifications   Permission = "notifications:view"
	PermissionCreateNotifications Permission = "notifications:create"
	PermissionDeleteNotifications Permission = "notifications:delete"

	// Report permissions
	PermissionViewReports   Permission = "reports:view"
	PermissionCreateReports Permission = "reports:create"
	PermissionUpdateReports Permission = "reports:update"
	PermissionDeleteReports Permission = "reports:delete"
	PermissionExportReports Permission = "reports:export"

	// Dashboard permissions
	PermissionViewDashboard Permission = "dashboard:view"

	// Settings permissions
	PermissionViewSettings   Permission = "settings:view"
	PermissionUpdateSettings Permission = "settings:update"

	// Task permissions
	PermissionViewTasks   Permission = "tasks:view"
	PermissionCreateTasks Permission = "tasks:create"
	PermissionUpdateTasks Permission = "tasks:update"
	PermissionDeleteTasks Permission = "tasks:delete"

	// Document permissions
	PermissionViewDocuments   Permission = "documents:view"
	PermissionCreateDocuments Permission = "documents:create"
	PermissionUpdateDocuments Permission = "documents:update"
	PermissionDeleteDocuments Permission = "documents:delete"

	// Partner permissions
	PermissionViewPartners   Permission = "partners:view"
	PermissionCreatePartners Permission = "partners:create"
	PermissionUpdatePartners Permission = "partners:update"
	PermissionDeletePartners Permission = "partners:delete"

	// Transfer permissions
	PermissionViewTransfers   Permission = "transfers:view"
	PermissionCreateTransfers Permission = "transfers:create"
	PermissionUpdateTransfers Permission = "transfers:update"
	PermissionDeleteTransfers Permission = "transfers:delete"

	// Contact permissions
	PermissionViewContacts   Permission = "contacts:view"
	PermissionCreateContacts Permission = "contacts:create"
	PermissionUpdateContacts Permission = "contacts:update"
	PermissionDeleteContacts Permission = "contacts:delete"

	// Inventory permissions
	PermissionViewInventory   Permission = "inventory:view"
	PermissionCreateInventory Permission = "inventory:create"
	PermissionUpdateInventory Permission = "inventory:update"
	PermissionDeleteInventory Permission = "inventory:delete"

	// Stock transaction permissions
	PermissionViewStockTransactions Permission = "stock:view"
)

// PermissionMatrix defines permissions for each role
var PermissionMatrix = map[entities.UserRole][]Permission{
	entities.RoleSuperAdmin: {
		// Super admin has all permissions
		PermissionViewUsers, PermissionCreateUsers, PermissionUpdateUsers, PermissionDeleteUsers,
		PermissionViewAnimals, PermissionCreateAnimals, PermissionUpdateAnimals, PermissionDeleteAnimals,
		PermissionViewAdoptions, PermissionCreateAdoptions, PermissionUpdateAdoptions, PermissionDeleteAdoptions,
		PermissionViewVeterinary, PermissionCreateVeterinary, PermissionUpdateVeterinary, PermissionDeleteVeterinary,
		PermissionViewDonors, PermissionCreateDonors, PermissionUpdateDonors, PermissionDeleteDonors,
		PermissionViewDonations, PermissionCreateDonations, PermissionUpdateDonations, PermissionDeleteDonations,
		PermissionViewContacts, PermissionCreateContacts, PermissionUpdateContacts, PermissionDeleteContacts,
		PermissionViewCampaigns, PermissionCreateCampaigns, PermissionUpdateCampaigns, PermissionDeleteCampaigns,
		PermissionViewEvents, PermissionCreateEvents, PermissionUpdateEvents, PermissionDeleteEvents,
		PermissionViewVolunteers, PermissionCreateVolunteers, PermissionUpdateVolunteers, PermissionDeleteVolunteers,
		PermissionViewCommunications, PermissionCreateCommunications, PermissionUpdateCommunications, PermissionDeleteCommunications,
		PermissionViewTemplates, PermissionCreateTemplates, PermissionUpdateTemplates, PermissionDeleteTemplates,
		PermissionViewNotifications, PermissionCreateNotifications, PermissionDeleteNotifications,
		PermissionViewReports, PermissionCreateReports, PermissionUpdateReports, PermissionDeleteReports, PermissionExportReports,
		PermissionViewDashboard,
		PermissionViewSettings, PermissionUpdateSettings,
		PermissionViewTasks, PermissionCreateTasks, PermissionUpdateTasks, PermissionDeleteTasks,
		PermissionViewDocuments, PermissionCreateDocuments, PermissionUpdateDocuments, PermissionDeleteDocuments,
		PermissionViewPartners, PermissionCreatePartners, PermissionUpdatePartners, PermissionDeletePartners,
		PermissionViewTransfers, PermissionCreateTransfers, PermissionUpdateTransfers, PermissionDeleteTransfers,
		PermissionViewInventory, PermissionCreateInventory, PermissionUpdateInventory, PermissionDeleteInventory,
		PermissionViewStockTransactions,
	},
	entities.RoleAdmin: {
		// Admin has most permissions except user management
		PermissionViewUsers,
		PermissionViewAnimals, PermissionCreateAnimals, PermissionUpdateAnimals, PermissionDeleteAnimals,
		PermissionViewAdoptions, PermissionCreateAdoptions, PermissionUpdateAdoptions, PermissionDeleteAdoptions,
		PermissionViewVeterinary, PermissionCreateVeterinary, PermissionUpdateVeterinary, PermissionDeleteVeterinary,
		PermissionViewDonors, PermissionCreateDonors, PermissionUpdateDonors, PermissionDeleteDonors,
		PermissionViewDonations, PermissionCreateDonations, PermissionUpdateDonations, PermissionDeleteDonations,
		PermissionViewContacts, PermissionCreateContacts, PermissionUpdateContacts, PermissionDeleteContacts,
		PermissionViewCampaigns, PermissionCreateCampaigns, PermissionUpdateCampaigns, PermissionDeleteCampaigns,
		PermissionViewEvents, PermissionCreateEvents, PermissionUpdateEvents, PermissionDeleteEvents,
		PermissionViewVolunteers, PermissionCreateVolunteers, PermissionUpdateVolunteers, PermissionDeleteVolunteers,
		PermissionViewCommunications, PermissionCreateCommunications, PermissionUpdateCommunications, PermissionDeleteCommunications,
		PermissionViewTemplates, PermissionCreateTemplates, PermissionUpdateTemplates, PermissionDeleteTemplates,
		PermissionViewNotifications, PermissionCreateNotifications, PermissionDeleteNotifications,
		PermissionViewReports, PermissionCreateReports, PermissionUpdateReports, PermissionDeleteReports, PermissionExportReports,
		PermissionViewDashboard,
		PermissionViewSettings, PermissionUpdateSettings,
		PermissionViewTasks, PermissionCreateTasks, PermissionUpdateTasks, PermissionDeleteTasks,
		PermissionViewDocuments, PermissionCreateDocuments, PermissionUpdateDocuments, PermissionDeleteDocuments,
		PermissionViewPartners, PermissionCreatePartners, PermissionUpdatePartners, PermissionDeletePartners,
		PermissionViewTransfers, PermissionCreateTransfers, PermissionUpdateTransfers, PermissionDeleteTransfers,
		PermissionViewInventory, PermissionCreateInventory, PermissionUpdateInventory, PermissionDeleteInventory,
		PermissionViewStockTransactions,
	},
	entities.RoleEmployee: {
		// Employee can manage animals, adoptions, and veterinary records
		PermissionViewAnimals, PermissionCreateAnimals, PermissionUpdateAnimals,
		PermissionViewAdoptions, PermissionCreateAdoptions, PermissionUpdateAdoptions,
		PermissionViewVeterinary, PermissionCreateVeterinary, PermissionUpdateVeterinary,
		PermissionViewDonors, PermissionCreateDonors,
		PermissionViewDonations, PermissionCreateDonations,
		PermissionViewCampaigns,
		PermissionViewContacts, PermissionCreateContacts, PermissionUpdateContacts,
		PermissionViewEvents, PermissionCreateEvents, PermissionUpdateEvents,
		PermissionViewVolunteers, PermissionCreateVolunteers, PermissionUpdateVolunteers,
		PermissionViewCommunications, PermissionCreateCommunications,
		PermissionViewTemplates,
		PermissionViewNotifications,
		PermissionViewReports,
		PermissionViewDashboard,
		PermissionViewTasks, PermissionCreateTasks, PermissionUpdateTasks,
		PermissionViewDocuments, PermissionCreateDocuments, PermissionUpdateDocuments,
		PermissionViewPartners, PermissionCreatePartners, PermissionUpdatePartners,
		PermissionViewTransfers, PermissionCreateTransfers, PermissionUpdateTransfers,
		PermissionViewInventory, PermissionCreateInventory, PermissionUpdateInventory,
		PermissionViewStockTransactions,
	},
	entities.RoleVolunteer: {
		// Volunteer has read access to most things, limited write access
		PermissionViewAnimals,
		PermissionViewAdoptions,
		PermissionViewVeterinary,
		PermissionViewDonors,
		PermissionViewDonations,
		PermissionViewCampaigns,
		PermissionViewContacts,
		PermissionViewEvents,
		PermissionViewVolunteers,
		PermissionViewNotifications,
		PermissionViewReports,
		PermissionViewTasks,
		PermissionViewDocuments,
		PermissionViewPartners,
		PermissionViewTransfers,
		PermissionViewInventory,
		PermissionViewStockTransactions,
	},
	entities.RoleUser: {
		// Regular user has minimal access
		PermissionViewAnimals,
		PermissionViewEvents,
		PermissionViewNotifications,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role entities.UserRole, permission Permission) bool {
	permissions, ok := PermissionMatrix[role]
	if !ok {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// RequirePermission middleware checks if the user has the required permission
func RequirePermission(permission Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		userEntity, ok := user.(*entities.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		if !HasPermission(userEntity.Role, permission) {
			c.JSON(http.StatusForbidden, errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware checks if the user has one of the required roles
func RequireRole(roles ...entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		userEntity, ok := user.(*entities.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if userEntity.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware checks if the user is an admin or super admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(entities.RoleSuperAdmin, entities.RoleAdmin)
}

// RequireSuperAdmin middleware checks if the user is a super admin
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(entities.RoleSuperAdmin)
}
