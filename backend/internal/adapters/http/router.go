package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/adapters/auth"
	"github.com/sainaif/animalsys/internal/adapters/http/handlers"
	"github.com/sainaif/animalsys/internal/infrastructure/middleware"
)

type Router struct {
	engine *gin.Engine

	// Handlers
	authHandler          *handlers.AuthHandler
	userHandler          *handlers.UserHandler
	animalHandler        *handlers.AnimalHandler
	adoptionHandler      *handlers.AdoptionHandler
	volunteerHandler     *handlers.VolunteerHandler
	scheduleHandler      *handlers.ScheduleHandler
	documentHandler      *handlers.DocumentHandler
	financeHandler       *handlers.FinanceHandler
	donorHandler         *handlers.DonorHandler
	inventoryHandler     *handlers.InventoryHandler
	veterinaryHandler    *handlers.VeterinaryHandler
	campaignHandler      *handlers.CampaignHandler
	partnerHandler       *handlers.PartnerHandler
	communicationHandler *handlers.CommunicationHandler

	// Middleware
	authMiddleware *auth.AuthMiddleware
	rbac           *auth.RBAC
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	animalHandler *handlers.AnimalHandler,
	adoptionHandler *handlers.AdoptionHandler,
	volunteerHandler *handlers.VolunteerHandler,
	scheduleHandler *handlers.ScheduleHandler,
	documentHandler *handlers.DocumentHandler,
	financeHandler *handlers.FinanceHandler,
	donorHandler *handlers.DonorHandler,
	inventoryHandler *handlers.InventoryHandler,
	veterinaryHandler *handlers.VeterinaryHandler,
	campaignHandler *handlers.CampaignHandler,
	partnerHandler *handlers.PartnerHandler,
	communicationHandler *handlers.CommunicationHandler,
	authMiddleware *auth.AuthMiddleware,
	rbac *auth.RBAC,
) *Router {
	engine := gin.Default()

	return &Router{
		engine:               engine,
		authHandler:          authHandler,
		userHandler:          userHandler,
		animalHandler:        animalHandler,
		adoptionHandler:      adoptionHandler,
		volunteerHandler:     volunteerHandler,
		scheduleHandler:      scheduleHandler,
		documentHandler:      documentHandler,
		financeHandler:       financeHandler,
		donorHandler:         donorHandler,
		inventoryHandler:     inventoryHandler,
		veterinaryHandler:    veterinaryHandler,
		campaignHandler:      campaignHandler,
		partnerHandler:       partnerHandler,
		communicationHandler: communicationHandler,
		authMiddleware:       authMiddleware,
		rbac:                 rbac,
	}
}

func (r *Router) SetupRoutes(corsMiddleware, securityMiddleware, loggerMiddleware, rateLimitMiddleware gin.HandlerFunc) {
	// Apply global middleware
	r.engine.Use(corsMiddleware)
	r.engine.Use(securityMiddleware)
	r.engine.Use(loggerMiddleware)
	r.engine.Use(rateLimitMiddleware)
	r.engine.Use(gin.Recovery())

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.engine.Group("/api/v1")

	// Public routes
	r.setupPublicRoutes(api)

	// Authenticated routes
	r.setupAuthenticatedRoutes(api)
}

func (r *Router) setupPublicRoutes(api *gin.RouterGroup) {
	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}

	// Public animal listings (for adopters)
	api.GET("/animals/available", r.animalHandler.GetAvailableForAdoption)
	api.GET("/animals/:id", r.animalHandler.GetByID)

	// Public campaign listings
	api.GET("/campaigns/active", r.campaignHandler.GetActive)
}

func (r *Router) setupAuthenticatedRoutes(api *gin.RouterGroup) {
	// Apply auth middleware
	authenticated := api.Group("")
	authenticated.Use(r.authMiddleware.Authenticate())

	// Auth profile
	authenticated.GET("/auth/profile", r.authHandler.GetProfile)
	authenticated.POST("/auth/change-password", r.authHandler.ChangePassword)

	// Users (Admin only)
	users := authenticated.Group("/users")
	users.Use(r.rbac.RequireRole("admin"))
	{
		users.GET("", r.userHandler.List)
		users.GET("/:id", r.userHandler.GetByID)
		users.PUT("/:id", r.userHandler.Update)
		users.DELETE("/:id", r.userHandler.Delete)
	}

	// Animals (Employee+)
	animals := authenticated.Group("/animals")
	animals.Use(r.rbac.RequireRole("employee"))
	{
		animals.POST("", r.animalHandler.Create)
		animals.GET("", r.animalHandler.List)
		animals.PUT("/:id", r.animalHandler.Update)
		animals.DELETE("/:id", r.animalHandler.Delete)
		animals.POST("/:id/medical-records", r.animalHandler.AddMedicalRecord)
		animals.POST("/:id/photos", r.animalHandler.AddPhoto)
	}

	// Adoptions
	adoptions := authenticated.Group("/adoptions")
	{
		adoptions.POST("", r.adoptionHandler.Create) // User can apply
		adoptions.GET("", r.adoptionHandler.List)
		adoptions.GET("/:id", r.adoptionHandler.GetByID)

		// Admin/Employee actions
		adoptionsAdmin := adoptions.Group("")
		adoptionsAdmin.Use(r.rbac.RequireRole("employee"))
		{
			adoptionsAdmin.PUT("/:id", r.adoptionHandler.Update)
			adoptionsAdmin.DELETE("/:id", r.adoptionHandler.Delete)
			adoptionsAdmin.POST("/:id/approve", r.adoptionHandler.Approve)
			adoptionsAdmin.POST("/:id/reject", r.adoptionHandler.Reject)
			adoptionsAdmin.POST("/:id/complete", r.adoptionHandler.Complete)
		}
	}

	// Volunteers
	volunteers := authenticated.Group("/volunteers")
	volunteers.Use(r.rbac.RequireRole("employee"))
	{
		volunteers.POST("", r.volunteerHandler.Create)
		volunteers.GET("", r.volunteerHandler.List)
		volunteers.GET("/active", r.volunteerHandler.GetActive)
		volunteers.GET("/:id", r.volunteerHandler.GetByID)
		volunteers.PUT("/:id", r.volunteerHandler.Update)
		volunteers.DELETE("/:id", r.volunteerHandler.Delete)
		volunteers.POST("/:id/trainings", r.volunteerHandler.AddTraining)
		volunteers.POST("/hours", r.volunteerHandler.LogHours)
		volunteers.GET("/:id/hours", r.volunteerHandler.GetVolunteerHours)
		volunteers.GET("/:id/statistics", r.volunteerHandler.GetStatistics)
	}

	// Schedules
	schedules := authenticated.Group("/schedules")
	schedules.Use(r.rbac.RequireRole("volunteer"))
	{
		schedules.POST("", r.scheduleHandler.Create)
		schedules.GET("", r.scheduleHandler.List)
		schedules.GET("/:id", r.scheduleHandler.GetByID)
		schedules.PUT("/:id", r.scheduleHandler.Update)
		schedules.DELETE("/:id", r.scheduleHandler.Delete)
		schedules.POST("/:id/request-swap", r.scheduleHandler.RequestSwap)
		schedules.POST("/:id/approve-swap", r.scheduleHandler.ApproveSwap)
		schedules.POST("/:id/reject-swap", r.scheduleHandler.RejectSwap)
		schedules.POST("/:id/complete", r.scheduleHandler.MarkComplete)
		schedules.POST("/:id/absent", r.scheduleHandler.MarkAbsent)
	}

	// Documents
	documents := authenticated.Group("/documents")
	documents.Use(r.rbac.RequireRole("employee"))
	{
		documents.POST("", r.documentHandler.Create)
		documents.GET("", r.documentHandler.List)
		documents.GET("/expiring", r.documentHandler.GetExpiringSoon)
		documents.GET("/:id", r.documentHandler.GetByID)
		documents.PUT("/:id", r.documentHandler.Update)
		documents.DELETE("/:id", r.documentHandler.Delete)
	}

	// Finance
	finance := authenticated.Group("/finance")
	finance.Use(r.rbac.RequireRole("employee"))
	{
		finance.POST("", r.financeHandler.Create)
		finance.GET("", r.financeHandler.List)
		finance.GET("/summary", r.financeHandler.GetSummary)
		finance.GET("/report", r.financeHandler.GetReport)
		finance.GET("/:id", r.financeHandler.GetByID)
		finance.PUT("/:id", r.financeHandler.Update)
		finance.DELETE("/:id", r.financeHandler.Delete)
	}

	// Donors
	donors := authenticated.Group("/donors")
	donors.Use(r.rbac.RequireRole("employee"))
	{
		donors.POST("", r.donorHandler.CreateDonor)
		donors.GET("", r.donorHandler.ListDonors)
		donors.GET("/top", r.donorHandler.GetTopDonors)
		donors.GET("/:id", r.donorHandler.GetDonorByID)
		donors.PUT("/:id", r.donorHandler.UpdateDonor)
		donors.DELETE("/:id", r.donorHandler.DeleteDonor)
	}

	// Donations
	donations := authenticated.Group("/donations")
	donations.Use(r.rbac.RequireRole("employee"))
	{
		donations.POST("", r.donorHandler.RecordDonation)
		donations.GET("", r.donorHandler.ListDonations)
		donations.GET("/statistics", r.donorHandler.GetDonationStatistics)
		donations.GET("/:id", r.donorHandler.GetDonationByID)
		donations.PUT("/:id", r.donorHandler.UpdateDonation)
		donations.DELETE("/:id", r.donorHandler.DeleteDonation)
	}

	// Inventory
	inventory := authenticated.Group("/inventory")
	inventory.Use(r.rbac.RequireRole("employee"))
	{
		inventory.POST("", r.inventoryHandler.CreateItem)
		inventory.GET("", r.inventoryHandler.ListItems)
		inventory.GET("/low-stock", r.inventoryHandler.GetLowStockItems)
		inventory.GET("/expiring", r.inventoryHandler.GetExpiringItems)
		inventory.GET("/statistics", r.inventoryHandler.GetStatistics)
		inventory.GET("/:id", r.inventoryHandler.GetItemByID)
		inventory.PUT("/:id", r.inventoryHandler.UpdateItem)
		inventory.DELETE("/:id", r.inventoryHandler.DeleteItem)
		inventory.POST("/:id/adjust", r.inventoryHandler.AdjustStock)
		inventory.GET("/:id/movements", r.inventoryHandler.GetItemStockMovements)
		inventory.POST("/movements", r.inventoryHandler.RecordStockMovement)
	}

	// Veterinary
	veterinary := authenticated.Group("/veterinary")
	veterinary.Use(r.rbac.RequireRole("employee"))
	{
		// Visits
		veterinary.POST("/visits", r.veterinaryHandler.RecordVisit)
		veterinary.GET("/visits", r.veterinaryHandler.ListVisits)
		veterinary.GET("/visits/upcoming-followups", r.veterinaryHandler.GetUpcomingFollowUps)
		veterinary.GET("/visits/statistics", r.veterinaryHandler.GetStatistics)
		veterinary.GET("/visits/:id", r.veterinaryHandler.GetVisitByID)
		veterinary.PUT("/visits/:id", r.veterinaryHandler.UpdateVisit)
		veterinary.DELETE("/visits/:id", r.veterinaryHandler.DeleteVisit)

		// Vaccinations
		veterinary.POST("/vaccinations", r.veterinaryHandler.RecordVaccination)
		veterinary.GET("/vaccinations", r.veterinaryHandler.ListVaccinations)
		veterinary.GET("/vaccinations/upcoming", r.veterinaryHandler.GetUpcomingVaccinations)
		veterinary.GET("/vaccinations/:id", r.veterinaryHandler.GetVaccinationByID)
		veterinary.PUT("/vaccinations/:id", r.veterinaryHandler.UpdateVaccination)
		veterinary.DELETE("/vaccinations/:id", r.veterinaryHandler.DeleteVaccination)
	}

	// Campaigns
	campaigns := authenticated.Group("/campaigns")
	campaigns.Use(r.rbac.RequireRole("employee"))
	{
		campaigns.POST("", r.campaignHandler.Create)
		campaigns.GET("", r.campaignHandler.List)
		campaigns.GET("/statistics", r.campaignHandler.GetAllStatistics)
		campaigns.GET("/:id", r.campaignHandler.GetByID)
		campaigns.GET("/:id/statistics", r.campaignHandler.GetStatistics)
		campaigns.PUT("/:id", r.campaignHandler.Update)
		campaigns.DELETE("/:id", r.campaignHandler.Delete)
		campaigns.POST("/:id/progress", r.campaignHandler.UpdateProgress)
	}

	// Partners
	partners := authenticated.Group("/partners")
	partners.Use(r.rbac.RequireRole("employee"))
	{
		partners.POST("", r.partnerHandler.Create)
		partners.GET("", r.partnerHandler.List)
		partners.GET("/active", r.partnerHandler.GetActive)
		partners.GET("/statistics", r.partnerHandler.GetStatistics)
		partners.GET("/:id", r.partnerHandler.GetByID)
		partners.PUT("/:id", r.partnerHandler.Update)
		partners.DELETE("/:id", r.partnerHandler.Delete)
		partners.POST("/:id/collaborations", r.partnerHandler.AddCollaboration)
	}

	// Communications
	communications := authenticated.Group("/communications")
	communications.Use(r.rbac.RequireRole("employee"))
	{
		communications.POST("", r.communicationHandler.CreateCommunication)
		communications.GET("", r.communicationHandler.ListCommunications)
		communications.GET("/scheduled", r.communicationHandler.GetScheduled)
		communications.GET("/statistics", r.communicationHandler.GetStatistics)
		communications.GET("/:id", r.communicationHandler.GetCommunicationByID)
		communications.PUT("/:id", r.communicationHandler.UpdateCommunication)
		communications.DELETE("/:id", r.communicationHandler.DeleteCommunication)
		communications.POST("/:id/send", r.communicationHandler.SendCommunication)

		// Templates
		communications.POST("/templates", r.communicationHandler.CreateTemplate)
		communications.GET("/templates", r.communicationHandler.ListTemplates)
		communications.GET("/templates/:id", r.communicationHandler.GetTemplateByID)
		communications.PUT("/templates/:id", r.communicationHandler.UpdateTemplate)
		communications.DELETE("/templates/:id", r.communicationHandler.DeleteTemplate)
		communications.POST("/templates/:id/create", r.communicationHandler.CreateFromTemplate)
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
