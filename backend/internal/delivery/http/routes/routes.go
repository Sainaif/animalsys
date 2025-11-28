package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/delivery/http/handlers"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/pkg/security"
)

// SetupRoutes sets up all application routes
func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	animalHandler *handlers.AnimalHandler,
	veterinaryHandler *handlers.VeterinaryHandler,
	adoptionHandler *handlers.AdoptionHandler,
	donorHandler *handlers.DonorHandler,
	donationHandler *handlers.DonationHandler,
	campaignHandler *handlers.CampaignHandler,
	eventHandler *handlers.EventHandler,
	volunteerHandler *handlers.VolunteerHandler,
	contactHandler *handlers.ContactHandler,
	communicationHandler *handlers.CommunicationHandler,
	notificationHandler *handlers.NotificationHandler,
	reportHandler *handlers.ReportHandler,
	dashboardHandler *handlers.DashboardHandler,
	settingsHandler *handlers.SettingsHandler,
	taskHandler *handlers.TaskHandler,
	documentHandler *handlers.DocumentHandler,
	partnerHandler *handlers.PartnerHandler,
	transferHandler *handlers.TransferHandler,
	inventoryHandler *handlers.InventoryHandler,
	stockTransactionHandler *handlers.StockTransactionHandler,
	auditLogHandler *handlers.AuditLogHandler,
	monitoringHandler *handlers.MonitoringHandler,
	medicalHandler *handlers.MedicalHandler,
	batchHandler *handlers.BatchHandler,
	jwtService *security.JWTService,
	userRepo repositories.UserRepository,
) {
	// Public routes (no authentication required)
	public := router.Group("/api/v1")
	{
		// Health check
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Auth routes (public)
		auth := public.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		public.POST("/public/donations", donationHandler.CreatePublicDonation)
	}

	// Protected routes (authentication required)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService, userRepo))
	{
		// Auth routes (protected)
		auth := protected.Group("/auth")
		{
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", authHandler.GetMe)
			auth.PUT("/change-password", authHandler.ChangePassword)

			// Register requires admin role
			auth.POST("/register",
				middleware.RequireAdmin(),
				authHandler.Register,
			)
		}

		// User management routes (admin only)
		users := protected.Group("/users")
		users.Use(middleware.RequirePermission(middleware.PermissionViewUsers))
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)

			// Create, update, delete require admin
			users.POST("",
				middleware.RequireAdmin(),
				userHandler.CreateUser,
			)

			users.PUT("/:id",
				middleware.RequireAdmin(),
				userHandler.UpdateUser,
			)

			users.DELETE("/:id",
				middleware.RequireSuperAdmin(),
				userHandler.DeleteUser,
			)

			users.PUT("/:id/reset-password",
				middleware.RequireAdmin(),
				userHandler.ResetPassword,
			)

			// Update user role
			users.PUT("/:id/role",
				middleware.RequireAdmin(),
				userHandler.UpdateUserRole,
			)

			// Update user status
			users.PUT("/:id/status",
				middleware.RequireAdmin(),
				userHandler.UpdateUserStatus,
			)
		}

		// Animal management routes
		animals := protected.Group("/animals")
		{
			// Public animal routes (view permission)
			animals.GET("",
				middleware.RequirePermission(middleware.PermissionViewAnimals),
				animalHandler.ListAnimals,
			)

			animals.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewAnimals),
				animalHandler.GetStatistics,
			)

			animals.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewAnimals),
				animalHandler.GetAnimal,
			)

			// Create animal (employees and above)
			animals.POST("",
				middleware.RequirePermission(middleware.PermissionCreateAnimals),
				animalHandler.CreateAnimal,
			)

			// Update animal (employees and above)
			animals.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateAnimals),
				animalHandler.UpdateAnimal,
			)

			// Delete animal (admin only)
			animals.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteAnimals),
				animalHandler.DeleteAnimal,
			)

			// Image upload (employees and above)
			animals.POST("/:id/images",
				middleware.RequirePermission(middleware.PermissionUpdateAnimals),
				animalHandler.UploadAnimalImages,
			)

			// Add daily note (employees and above)
			animals.POST("/:id/notes",
				middleware.RequirePermission(middleware.PermissionUpdateAnimals),
				animalHandler.AddDailyNote,
			)

			// Animal-specific veterinary routes
			animals.GET("/:id/visits",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				veterinaryHandler.GetVisitsByAnimal,
			)

			animals.GET("/:id/vaccinations",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				veterinaryHandler.GetVaccinationsByAnimal,
			)
		}

		// Veterinary management routes
		veterinary := protected.Group("/veterinary")
		{
			// Veterinary visit routes
			visits := veterinary.Group("/visits")
			{
				// List and view visits
				visits.GET("",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.ListVisits,
				)

				visits.GET("/upcoming",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.GetUpcomingVisits,
				)

				visits.GET("/:id",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.GetVisit,
				)

				// Create visit (employees and above)
				visits.POST("",
					middleware.RequirePermission(middleware.PermissionCreateVeterinary),
					veterinaryHandler.CreateVisit,
				)

				// Update visit (employees and above)
				visits.PUT("/:id",
					middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
					veterinaryHandler.UpdateVisit,
				)

				// Delete visit (admin only)
				visits.DELETE("/:id",
					middleware.RequirePermission(middleware.PermissionDeleteVeterinary),
					veterinaryHandler.DeleteVisit,
				)
			}

			// Vaccination routes
			vaccinations := veterinary.Group("/vaccinations")
			{
				// List and view vaccinations
				vaccinations.GET("",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.ListVaccinations,
				)

				vaccinations.GET("/due",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.GetDueVaccinations,
				)

				vaccinations.GET("/:id",
					middleware.RequirePermission(middleware.PermissionViewVeterinary),
					veterinaryHandler.GetVaccination,
				)

				// Create vaccination (employees and above)
				vaccinations.POST("",
					middleware.RequirePermission(middleware.PermissionCreateVeterinary),
					veterinaryHandler.CreateVaccination,
				)

				// Delete vaccination (admin only)
				vaccinations.DELETE("/:id",
					middleware.RequirePermission(middleware.PermissionDeleteVeterinary),
					veterinaryHandler.DeleteVaccination,
				)
			}

			// Veterinary records routes
			veterinary.GET("/records",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				veterinaryHandler.ListVeterinaryRecords,
			)

			veterinary.POST("/records",
				middleware.RequirePermission(middleware.PermissionCreateVeterinary),
				veterinaryHandler.CreateVeterinaryRecord,
			)
		}

		// Adoption management routes
		adoptions := protected.Group("/adoptions")
		{
			// Adoption application routes
			applications := adoptions.Group("/applications")
			{
				// List and view applications
				applications.GET("",
					middleware.RequirePermission(middleware.PermissionViewAdoptions),
					adoptionHandler.ListApplications,
				)

				applications.GET("/pending",
					middleware.RequirePermission(middleware.PermissionViewAdoptions),
					adoptionHandler.GetPendingApplications,
				)

				applications.GET("/:id",
					middleware.RequirePermission(middleware.PermissionViewAdoptions),
					adoptionHandler.GetApplication,
				)

				// Create application (anyone can apply)
				applications.POST("",
					adoptionHandler.CreateApplication,
				)

				// Update application (employees and above - for review/approval)
				applications.PUT("/:id",
					middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
					adoptionHandler.UpdateApplication,
				)

				// Delete application (admin only)
				applications.DELETE("/:id",
					middleware.RequirePermission(middleware.PermissionDeleteAdoptions),
					adoptionHandler.DeleteApplication,
				)

				// Schedule visit for application
				applications.POST("/:id/schedule-visit",
					middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
					adoptionHandler.ScheduleVisit,
				)

				// Record home visit
				applications.POST("/:id/record-home-visit",
					middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
					adoptionHandler.RecordHomeVisit,
				)

				// Get visits for application
				applications.GET("/:id/visits",
					middleware.RequirePermission(middleware.PermissionViewAdoptions),
					adoptionHandler.GetVisits,
				)

				// Approve application
				applications.POST("/:id/approve",
					middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
					adoptionHandler.ApproveApplication,
				)

				// Reject application
				applications.POST("/:id/reject",
					middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
					adoptionHandler.RejectApplication,
				)
			}

			// Adoption routes
			adoptions.GET("",
				middleware.RequirePermission(middleware.PermissionViewAdoptions),
				adoptionHandler.ListAdoptions,
			)

			adoptions.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewAdoptions),
				adoptionHandler.GetAdoptionStatistics,
			)

			adoptions.GET("/follow-ups/pending",
				middleware.RequirePermission(middleware.PermissionViewAdoptions),
				adoptionHandler.GetPendingFollowUps,
			)

			adoptions.GET("/contracts",
				middleware.RequirePermission(middleware.PermissionViewAdoptions),
				adoptionHandler.GetContracts,
			)

			adoptions.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewAdoptions),
				adoptionHandler.GetAdoption,
			)

			// Create adoption (employees and above)
			adoptions.POST("",
				middleware.RequirePermission(middleware.PermissionCreateAdoptions),
				adoptionHandler.CreateAdoption,
			)

			// Update adoption (employees and above)
			adoptions.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
				adoptionHandler.UpdateAdoption,
			)

			// Delete adoption (admin only)
			adoptions.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteAdoptions),
				adoptionHandler.DeleteAdoption,
			)

			// Finalize adoption
			adoptions.POST("/:id/finalize",
				middleware.RequirePermission(middleware.PermissionUpdateAdoptions),
				adoptionHandler.FinalizeAdoption,
			)
		}

		// Animal-specific adoption routes
		animals.GET("/:id/applications",
			middleware.RequirePermission(middleware.PermissionViewAdoptions),
			adoptionHandler.GetApplicationsByAnimal,
		)

		animals.GET("/:id/adoption",
			middleware.RequirePermission(middleware.PermissionViewAdoptions),
			adoptionHandler.GetAdoptionByAnimal,
		)

		// Donor management routes
		donors := protected.Group("/donors")
		{
			// List and view donors
			donors.GET("",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.ListDonors,
			)

			donors.GET("/major",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetMajorDonors,
			)

			donors.GET("/lapsed",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetLapsedDonors,
			)

			donors.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetDonorStatistics,
			)

			donors.GET("/top",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetTopDonors,
			)

			donors.GET("/recurring",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetRecurringDonors,
			)

			donors.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewDonors),
				donorHandler.GetDonor,
			)

			// Create donor (employees and above)
			donors.POST("",
				middleware.RequirePermission(middleware.PermissionCreateDonors),
				donorHandler.CreateDonor,
			)

			// Update donor (employees and above)
			donors.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateDonors),
				donorHandler.UpdateDonor,
			)

			// Delete donor (admin only)
			donors.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteDonors),
				donorHandler.DeleteDonor,
			)

			// Update engagement (employees and above)
			donors.POST("/:id/engagement",
				middleware.RequirePermission(middleware.PermissionUpdateDonors),
				donorHandler.UpdateDonorEngagement,
			)

			// Update communication preferences
			donors.PUT("/:id/communication-preferences",
				middleware.RequirePermission(middleware.PermissionUpdateDonors),
				donorHandler.UpdateCommunicationPreferences,
			)
		}

		// Donation management routes
		donations := protected.Group("/donations")
		{
			// List and view donations
			donations.GET("",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.ListDonations,
			)

			donations.GET("/recurring",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetRecurringDonations,
			)

			donations.GET("/pending-thank-yous",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetPendingThankYous,
			)

			donations.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetDonationStatistics,
			)

			donations.GET("/date-range",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetDonationsByDateRange,
			)

			donations.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetDonation,
			)

			// Create donation (employees and above)
			donations.POST("",
				middleware.RequirePermission(middleware.PermissionCreateDonations),
				donationHandler.CreateDonation,
			)

			// Update donation (employees and above)
			donations.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateDonations),
				donationHandler.UpdateDonation,
			)

			// Delete donation (admin only)
			donations.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteDonations),
				donationHandler.DeleteDonation,
			)

			// Process donation (employees and above)
			donations.POST("/:id/process",
				middleware.RequirePermission(middleware.PermissionUpdateDonations),
				donationHandler.ProcessDonation,
			)

			// Refund donation (admin only)
			donations.POST("/:id/refund",
				middleware.RequirePermission(middleware.PermissionDeleteDonations),
				donationHandler.RefundDonation,
			)

			// Send thank you (employees and above)
			donations.POST("/:id/thank-you",
				middleware.RequirePermission(middleware.PermissionUpdateDonations),
				donationHandler.SendThankYou,
			)

			// Generate tax receipt (employees and above)
			donations.POST("/:id/tax-receipt",
				middleware.RequirePermission(middleware.PermissionUpdateDonations),
				donationHandler.GenerateTaxReceipt,
			)

			// Get donation receipt
			donations.GET("/:id/receipt",
				middleware.RequirePermission(middleware.PermissionViewDonations),
				donationHandler.GetDonationReceipt,
			)
		}

		// Campaign management routes
		campaigns := protected.Group("/campaigns")
		{
			// List and view campaigns
			campaigns.GET("",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.ListCampaigns,
			)

			campaigns.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetActiveCampaigns,
			)

			campaigns.GET("/featured",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetFeaturedCampaigns,
			)

			campaigns.GET("/public",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetPublicCampaigns,
			)

			campaigns.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetCampaignStatistics,
			)

			campaigns.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetCampaign,
			)

			// Create campaign (admin only)
			campaigns.POST("",
				middleware.RequirePermission(middleware.PermissionCreateCampaigns),
				campaignHandler.CreateCampaign,
			)

			// Update campaign (admin only)
			campaigns.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.UpdateCampaign,
			)

			// Delete campaign (admin only)
			campaigns.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteCampaigns),
				campaignHandler.DeleteCampaign,
			)

			// Campaign status actions (admin only)
			campaigns.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.ActivateCampaign,
			)

			campaigns.POST("/:id/pause",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.PauseCampaign,
			)

			campaigns.POST("/:id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.CompleteCampaign,
			)

			campaigns.POST("/:id/cancel",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.CancelCampaign,
			)

			// Get campaign-specific statistics
			campaigns.GET("/:id/statistics",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetCampaignStatistics,
			)

			// Get campaign donors
			campaigns.GET("/:id/donors",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetCampaignDonors,
			)

			// Update campaign amount
			campaigns.PUT("/:id/amount",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.UpdateCampaignAmount,
			)

			// Get campaign progress
			campaigns.GET("/:id/progress",
				middleware.RequirePermission(middleware.PermissionViewCampaigns),
				campaignHandler.GetCampaignProgress,
			)

			// Share campaign
			campaigns.POST("/:id/share",
				middleware.RequirePermission(middleware.PermissionUpdateCampaigns),
				campaignHandler.ShareCampaign,
			)
		}

		// Donor-specific routes (outside donors group to avoid route conflicts)
		protected.GET("/donors/:id/donations",
			middleware.RequirePermission(middleware.PermissionViewDonations),
			donationHandler.GetDonationsByDonor,
		)

		// Campaign-specific routes (outside campaigns group to avoid route conflicts)
		protected.GET("/campaigns/:id/donations",
			middleware.RequirePermission(middleware.PermissionViewDonations),
			donationHandler.GetDonationsByCampaign,
		)

		// User-specific routes (outside users group to avoid route conflicts)
		protected.GET("/users/:id/campaigns",
			middleware.RequirePermission(middleware.PermissionViewCampaigns),
			campaignHandler.GetCampaignsByManager,
		)

		// Event management routes
		events := protected.Group("/events")
		{
			// List and view events
			events.GET("",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.ListEvents,
			)

			events.GET("/upcoming",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetUpcomingEvents,
			)

			events.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetActiveEvents,
			)

			events.GET("/public",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetPublicEvents,
			)

			events.GET("/featured",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetFeaturedEvents,
			)

			events.GET("/needing-volunteers",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEventsNeedingVolunteers,
			)

			events.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEventStatistics,
			)

			events.GET("/past",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetPastEvents,
			)

			events.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEvent,
			)

			// Create event (employees and above)
			events.POST("",
				middleware.RequirePermission(middleware.PermissionCreateEvents),
				eventHandler.CreateEvent,
			)

			// Update event (employees and above)
			events.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.UpdateEvent,
			)

			// Delete event (admin only)
			events.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteEvents),
				eventHandler.DeleteEvent,
			)

			// Event status actions (employees and above)
			events.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.ActivateEvent,
			)

			events.POST("/:id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.CompleteEvent,
			)

			events.POST("/:id/cancel",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.CancelEvent,
			)

			// Volunteer assignment (employees and above)
			events.POST("/:id/assign-volunteer",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.AssignVolunteer,
			)

			events.POST("/:id/unassign-volunteer",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.UnassignVolunteer,
			)

			// Register for event
			events.POST("/:id/register",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.RegisterForEvent,
			)

			// Get event registrations
			events.GET("/:id/registrations",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEventRegistrations,
			)

			// Get event-specific statistics
			events.GET("/:id/statistics",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEventStatisticsDetail,
			)

			// Publish event
			events.POST("/:id/publish",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.PublishEvent,
			)

			// Send event reminder
			events.POST("/:id/send-reminder",
				middleware.RequirePermission(middleware.PermissionUpdateEvents),
				eventHandler.SendEventReminder,
			)

			// Get event attendance
			events.GET("/:id/attendance",
				middleware.RequirePermission(middleware.PermissionViewEvents),
				eventHandler.GetEventAttendance,
			)
		}

		// Volunteer management routes
		volunteers := protected.Group("/volunteers")
		{
			// List and view volunteers
			volunteers.GET("",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.ListVolunteers,
			)

			volunteers.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetActiveVolunteers,
			)

			volunteers.GET("/by-skill",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteersBySkill,
			)

			volunteers.GET("/needing-background-check",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteersNeedingBackgroundCheck,
			)

			volunteers.GET("/top",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetTopVolunteers,
			)

			volunteers.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteerStatistics,
			)

			volunteers.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteer,
			)

			// Create volunteer (employees and above)
			volunteers.POST("",
				middleware.RequirePermission(middleware.PermissionCreateVolunteers),
				volunteerHandler.CreateVolunteer,
			)

			// Update volunteer (employees and above)
			volunteers.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.UpdateVolunteer,
			)

			// Delete volunteer (admin only)
			volunteers.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteVolunteers),
				volunteerHandler.DeleteVolunteer,
			)

			// Volunteer management actions (employees and above)
			volunteers.POST("/:id/approve",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.ApproveVolunteer,
			)

			volunteers.POST("/:id/suspend",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.SuspendVolunteer,
			)


			volunteers.POST("/:id/commendation",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.AddCommendation,
			)

			volunteers.POST("/:id/warning",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.AddWarning,
			)

			volunteers.POST("/:id/certification",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.AddCertification,
			)

			// Log volunteer hours
			volunteers.POST("/:id/log-hours",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.LogVolunteerHours,
			)

			// Get volunteer hours
			volunteers.GET("/:id/hours",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteerHours,
			)

			// Activate volunteer
			volunteers.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.ActivateVolunteer,
			)

			// Deactivate volunteer
			volunteers.POST("/:id/deactivate",
				middleware.RequirePermission(middleware.PermissionUpdateVolunteers),
				volunteerHandler.DeactivateVolunteer,
			)

			// Get volunteers by role
			volunteers.GET("/role/:role",
				middleware.RequirePermission(middleware.PermissionViewVolunteers),
				volunteerHandler.GetVolunteersByRole,
			)
		}

		// Contact management routes
		contacts := protected.Group("/contacts")
		contacts.Use(middleware.RequirePermission(middleware.PermissionViewContacts))
		{
			contacts.GET("",
				contactHandler.List,
			)

			contacts.GET("/:id",
				contactHandler.Get,
			)

			contacts.POST("",
				middleware.RequirePermission(middleware.PermissionCreateContacts),
				contactHandler.Create,
			)

			contacts.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateContacts),
				contactHandler.Update,
			)

			contacts.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteContacts),
				contactHandler.Delete,
			)
		}

		// Communication and Template management routes
		communications := protected.Group("/communications")
		{
			// List and view communications
			communications.GET("",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.ListCommunications,
			)

			communications.GET("/pending",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetPendingCommunications,
			)

			communications.GET("/retry",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetCommunicationsForRetry,
			)

			communications.GET("/by-recipient",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetCommunicationsByRecipient,
			)

			communications.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetCommunicationStatistics,
			)

			communications.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetCommunication,
			)

			// Create communication (employees and above)
			communications.POST("",
				middleware.RequirePermission(middleware.PermissionCreateCommunications),
				communicationHandler.CreateCommunication,
			)

			communications.POST("/send-from-template",
				middleware.RequirePermission(middleware.PermissionCreateCommunications),
				communicationHandler.SendFromTemplate,
			)

			// Update communication status (employees and above)
			communications.PUT("/:id/status",
				middleware.RequirePermission(middleware.PermissionUpdateCommunications),
				communicationHandler.UpdateCommunicationStatus,
			)

			// Update communication
			communications.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateCommunications),
				communicationHandler.UpdateCommunication,
			)

			// Get communication status
			communications.GET("/:id/status",
				middleware.RequirePermission(middleware.PermissionViewCommunications),
				communicationHandler.GetCommunicationStatus,
			)

			// Delete communication (admin only)
			communications.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteCommunications),
				communicationHandler.DeleteCommunication,
			)

			// Tracking endpoints (no auth required for webhook-style tracking)
		}

		// Public tracking endpoints for email opens and clicks
		protected.POST("/communications/:id/track-open", communicationHandler.TrackOpen)
		protected.POST("/communications/:id/track-click", communicationHandler.TrackClick)

		// Campaign-specific communications
		protected.GET("/campaigns/:id/communications",
			middleware.RequirePermission(middleware.PermissionViewCommunications),
			communicationHandler.GetCommunicationsByCampaign,
		)

		// Batch communications
		protected.GET("/batches/:id/communications",
			middleware.RequirePermission(middleware.PermissionViewCommunications),
			communicationHandler.GetCommunicationsByBatch,
		)

		// Template management routes
		templates := protected.Group("/templates")
		{
			// List and view templates
			templates.GET("",
				middleware.RequirePermission(middleware.PermissionViewTemplates),
				communicationHandler.ListTemplates,
			)

			templates.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewTemplates),
				communicationHandler.GetActiveTemplates,
			)

			templates.GET("/by-category",
				middleware.RequirePermission(middleware.PermissionViewTemplates),
				communicationHandler.GetTemplatesByCategory,
			)

			templates.GET("/default",
				middleware.RequirePermission(middleware.PermissionViewTemplates),
				communicationHandler.GetDefaultTemplate,
			)

			templates.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewTemplates),
				communicationHandler.GetTemplate,
			)

			// Create template (admin only)
			templates.POST("",
				middleware.RequirePermission(middleware.PermissionCreateTemplates),
				communicationHandler.CreateTemplate,
			)

			// Update template (admin only)
			templates.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateTemplates),
				communicationHandler.UpdateTemplate,
			)

			// Delete template (admin only)
			templates.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteTemplates),
				communicationHandler.DeleteTemplate,
			)
		}

		// Notification management routes
		notifications := protected.Group("/notifications")
		{
			// User's own notifications (all authenticated users)
			notifications.GET("/me",
				notificationHandler.GetMyNotifications,
			)

			notifications.GET("/unread",
				notificationHandler.GetUnreadNotifications,
			)

			notifications.GET("/unread/count",
				notificationHandler.GetUnreadCount,
			)

			notifications.GET("/preferences",
				notificationHandler.GetNotificationPreferences,
			)

			// Update notification preferences
			notifications.PUT("/preferences",
				notificationHandler.UpdateNotificationPreferences,
			)

			// Get notifications by type
			notifications.GET("/type/:type",
				notificationHandler.GetNotificationsByType,
			)

			// List notifications with filters
			notifications.GET("",
				middleware.RequirePermission(middleware.PermissionViewNotifications),
				notificationHandler.ListNotifications,
			)

			notifications.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewNotifications),
				notificationHandler.GetNotification,
			)

			// Create notification (admin only)
			notifications.POST("",
				middleware.RequirePermission(middleware.PermissionCreateNotifications),
				notificationHandler.CreateNotification,
			)

			// Mark as read/unread (users can manage their own notifications)
			notifications.POST("/:id/read",
				notificationHandler.MarkAsRead,
			)

			notifications.POST("/:id/unread",
				notificationHandler.MarkAsUnread,
			)

			notifications.POST("/read-all",
				notificationHandler.MarkAllAsRead,
			)

			// Dismiss notification
			notifications.POST("/:id/dismiss",
				notificationHandler.DismissNotification,
			)

			// Delete notification
			notifications.DELETE("/:id",
				notificationHandler.DeleteNotification,
			)
		}

		// Report management routes
		reports := protected.Group("/reports")
		{
			// List and view reports
			reports.GET("",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.ListReports,
			)

			reports.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetActiveReports,
			)

			reports.GET("/public",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetPublicReports,
			)

			reports.GET("/executions/recent",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetRecentExecutions,
			)

			// Analytics report endpoints
			reports.GET("/animals",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetAnimalsReport,
			)

			reports.GET("/adoptions",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetAdoptionsReport,
			)

			reports.GET("/donations",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetDonationsReport,
			)

			reports.GET("/financials",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetFinancialsReport,
			)

			reports.GET("/volunteers",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetVolunteersReport,
			)

			reports.GET("/inventory",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetInventoryReport,
			)

			reports.GET("/veterinary",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetVeterinaryReport,
			)

			reports.GET("/compliance",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetComplianceReport,
			)

			reports.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetReport,
			)

			// Create report (admin only)
			reports.POST("",
				middleware.RequirePermission(middleware.PermissionCreateReports),
				reportHandler.CreateReport,
			)

			// Update report (admin only)
			reports.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateReports),
				reportHandler.UpdateReport,
			)

			// Delete report (admin only)
			reports.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteReports),
				reportHandler.DeleteReport,
			)

			// Execute report (view permission required)
			reports.POST("/:id/execute",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.ExecuteReport,
			)

			// Get report executions
			reports.GET("/:id/executions",
				middleware.RequirePermission(middleware.PermissionViewReports),
				reportHandler.GetReportExecutions,
			)

			// Create custom report
			reports.POST("/custom",
				middleware.RequirePermission(middleware.PermissionCreateReports),
				reportHandler.CreateCustomReport,
			)
		}

		// Dashboard routes
		dashboard := protected.Group("/dashboard")
		{
			// Get full dashboard metrics
			dashboard.GET("",
				middleware.RequirePermission(middleware.PermissionViewDashboard),
				dashboardHandler.GetDashboardMetrics,
			)

			// Get overview metrics only
			dashboard.GET("/overview",
				middleware.RequirePermission(middleware.PermissionViewDashboard),
				dashboardHandler.GetOverviewMetrics,
			)

			// Get dashboard summary
			dashboard.GET("/summary",
				middleware.RequirePermission(middleware.PermissionViewDashboard),
				dashboardHandler.GetDashboardSummary,
			)

			// Get dashboard widgets
			dashboard.GET("/widgets",
				middleware.RequirePermission(middleware.PermissionViewDashboard),
				dashboardHandler.GetDashboardWidgets,
			)
		}

		// Settings routes
		settings := protected.Group("/settings")
		{
			// Get settings
			settings.GET("",
				middleware.RequirePermission(middleware.PermissionViewSettings),
				settingsHandler.GetSettings,
			)

			// Update settings (admin only)
			settings.PUT("",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateSettings,
			)

			// Initialize settings (super admin only)
			settings.POST("/initialize",
				middleware.RequireSuperAdmin(),
				settingsHandler.InitializeSettings,
			)

			// Update specific sections (admin only)
			settings.PUT("/email",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateEmailSettings,
			)

			settings.PUT("/notifications",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateNotificationSettings,
			)

			settings.PUT("/features",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateFeatureFlags,
			)

			settings.PUT("/branding",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateBranding,
			)

			// Get organization settings
			settings.GET("/organization",
				middleware.RequirePermission(middleware.PermissionViewSettings),
				settingsHandler.GetOrganizationSettings,
			)

			// Update organization settings
			settings.PUT("/organization",
				middleware.RequirePermission(middleware.PermissionUpdateSettings),
				settingsHandler.UpdateOrganizationSettings,
			)

			// Get email settings
			settings.GET("/email",
				middleware.RequirePermission(middleware.PermissionViewSettings),
				settingsHandler.GetEmailSettings,
			)

			// Get notification settings
			settings.GET("/notifications",
				middleware.RequirePermission(middleware.PermissionViewSettings),
				settingsHandler.GetNotificationSettings,
			)

			// Get integration settings
			settings.GET("/integrations",
				middleware.RequirePermission(middleware.PermissionViewSettings),
				settingsHandler.GetIntegrationSettings,
			)

			// Public-accessible endpoints
			settings.GET("/contact",
				settingsHandler.GetContactInfo,
			)

			settings.GET("/hours",
				settingsHandler.GetOperatingHours,
			)
		}

		// Task management routes
		tasks := protected.Group("/tasks")
		{
			// List tasks
			tasks.GET("",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.ListTasks,
			)

			// Get my tasks
			tasks.GET("/my",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetMyTasks,
			)

			// Get overdue tasks
			tasks.GET("/overdue",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetOverdueTasks,
			)

			// Get upcoming tasks
			tasks.GET("/upcoming",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetUpcomingTasks,
			)

			// Get task statistics
			tasks.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetTaskStatistics,
			)

			// Get tasks by assignee
			tasks.GET("/assignee/:user_id",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetTasksByAssignee,
			)

			// Get task by ID
			tasks.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetTask,
			)

			// Create task
			tasks.POST("",
				middleware.RequirePermission(middleware.PermissionCreateTasks),
				taskHandler.CreateTask,
			)

			// Update task
			tasks.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.UpdateTask,
			)

			// Delete task
			tasks.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteTasks),
				taskHandler.DeleteTask,
			)

			// Task actions
			tasks.POST("/:id/assign",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.AssignTask,
			)

			tasks.POST("/:id/unassign",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.UnassignTask,
			)

			tasks.POST("/:id/start",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.StartTask,
			)

			tasks.POST("/:id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.CompleteTask,
			)

			tasks.POST("/:id/cancel",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.CancelTask,
			)

			// Checklist management
			tasks.POST("/:id/checklist",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.AddChecklistItem,
			)

			tasks.POST("/:id/checklist/:item_id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.CompleteChecklistItem,
			)

			tasks.DELETE("/:id/checklist/:item_id",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.RemoveChecklistItem,
			)

			// Add task comment
			tasks.POST("/:id/comments",
				middleware.RequirePermission(middleware.PermissionUpdateTasks),
				taskHandler.AddTaskComment,
			)

			// Get task comments
			tasks.GET("/:id/comments",
				middleware.RequirePermission(middleware.PermissionViewTasks),
				taskHandler.GetTaskComments,
			)
		}

		// Document management routes
		documents := protected.Group("/documents")
		{
			// List documents
			documents.GET("",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.ListDocuments,
			)

			// Get public documents
			documents.GET("/public",
				documentHandler.GetPublicDocuments,
			)

			// Get my documents
			documents.GET("/my",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetMyDocuments,
			)

			// Get expired documents
			documents.GET("/expired",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetExpiredDocuments,
			)

			// Get expiring soon documents
			documents.GET("/expiring-soon",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetExpiringSoonDocuments,
			)

			// Get document statistics
			documents.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocumentStatistics,
			)

			// Get document by ID
			documents.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocument,
			)

			// Get document versions
			documents.GET("/:id/versions",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocumentVersions,
			)

			// Download document
			documents.GET("/:id/download",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.DownloadDocument,
			)

			// Create document
			documents.POST("",
				middleware.RequirePermission(middleware.PermissionCreateDocuments),
				documentHandler.CreateDocument,
			)

			// Create new version
			documents.POST("/:id/versions",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.CreateNewVersion,
			)

			// Update document
			documents.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.UpdateDocument,
			)

			// Delete document
			documents.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteDocuments),
				documentHandler.DeleteDocument,
			)

			// Document access control
			documents.POST("/:id/grant-access",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.GrantAccess,
			)

			documents.POST("/:id/revoke-access",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.RevokeAccess,
			)

			documents.POST("/:id/make-public",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.MakePublic,
			)

			documents.POST("/:id/make-private",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.MakePrivate,
			)

			// Set expiration
			documents.POST("/:id/set-expiration",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.SetExpiration,
			)

			// Get documents by entity
			documents.GET("/entity/:entity_type/:entity_id",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocumentsByEntity,
			)

			// Get documents by type
			documents.GET("/type/:type",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocumentsByType,
			)

			// Get documents by category
			documents.GET("/category/:category",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.GetDocumentsByCategory,
			)

			// Search documents
			documents.GET("/search",
				middleware.RequirePermission(middleware.PermissionViewDocuments),
				documentHandler.SearchDocuments,
			)

			// Share document
			documents.POST("/:id/share",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.ShareDocument,
			)

			// Archive document
			documents.POST("/:id/archive",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.ArchiveDocument,
			)

			// Unarchive document
			documents.POST("/:id/unarchive",
				middleware.RequirePermission(middleware.PermissionUpdateDocuments),
				documentHandler.UnarchiveDocument,
			)
		}

		// Partner management routes
		partners := protected.Group("/partners")
		{
			// List partners
			partners.GET("",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.ListPartners,
			)

			// Get active partners
			partners.GET("/active",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetActivePartners,
			)

			// Get partners with capacity
			partners.GET("/with-capacity",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnersWithCapacity,
			)

			// Get partners by type
			partners.GET("/type/:type",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnersByType,
			)

			// Get partners by status
			partners.GET("/status/:status",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnersByStatus,
			)

			// Get partner statistics
			partners.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnerStatistics,
			)

			// Get partners accepting intakes
			partners.GET("/accepting-intakes",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnersAcceptingIntakes,
			)

			// Rate partner
			partners.POST("/:id/rate",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.RatePartner,
			)

			// Get partner-specific statistics
			partners.GET("/:id/statistics",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartnerStatisticsDetail,
			)

			// Update partner capacity
			partners.POST("/:id/update-capacity",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.UpdatePartnerCapacity,
			)

			// Get partner by ID
			partners.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewPartners),
				partnerHandler.GetPartner,
			)

			// Create partner
			partners.POST("",
				middleware.RequirePermission(middleware.PermissionCreatePartners),
				partnerHandler.CreatePartner,
			)

			// Update partner
			partners.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.UpdatePartner,
			)

			// Delete partner
			partners.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeletePartners),
				partnerHandler.DeletePartner,
			)

			// Partner actions
			partners.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.ActivatePartner,
			)

			partners.POST("/:id/suspend",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.SuspendPartner,
			)

			partners.POST("/:id/deactivate",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.DeactivatePartner,
			)

			partners.POST("/:id/add-rating",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.AddRating,
			)


			// Update partner capacity
			partners.PUT("/:id/capacity",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.UpdatePartnerCapacity,
			)

			partners.POST("/:id/set-agreement-expiry",
				middleware.RequirePermission(middleware.PermissionUpdatePartners),
				partnerHandler.SetAgreementExpiry,
			)
		}

		// Transfer management routes
		transfers := protected.Group("/transfers")
		{
			// List transfers
			transfers.GET("",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.ListTransfers,
			)

			// Get pending transfers
			transfers.GET("/pending",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetPendingTransfers,
			)

			// Get upcoming transfers
			transfers.GET("/upcoming",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetUpcomingTransfers,
			)

			// Get overdue transfers
			transfers.GET("/overdue",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetOverdueTransfers,
			)

			// Get transfers requiring follow-up
			transfers.GET("/follow-up",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetRequiringFollowUp,
			)

			// Get transfers by animal
			transfers.GET("/animal/:animal_id",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetTransfersByAnimal,
			)

			// Get transfers by partner
			transfers.GET("/partner/:partner_id",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetTransfersByPartner,
			)

			// Get transfers by status
			transfers.GET("/status/:status",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetTransfersByStatus,
			)

			// Get transfer statistics
			transfers.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetTransferStatistics,
			)

			// Get transfer by ID
			transfers.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewTransfers),
				transferHandler.GetTransfer,
			)

			// Create transfer
			transfers.POST("",
				middleware.RequirePermission(middleware.PermissionCreateTransfers),
				transferHandler.CreateTransfer,
			)

			// Update transfer
			transfers.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.UpdateTransfer,
			)

			// Delete transfer
			transfers.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteTransfers),
				transferHandler.DeleteTransfer,
			)

			// Transfer actions
			transfers.POST("/:id/approve",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.ApproveTransfer,
			)

			transfers.POST("/:id/reject",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.RejectTransfer,
			)

			transfers.POST("/:id/start-transit",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.StartTransit,
			)

			transfers.POST("/:id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.CompleteTransfer,
			)

			transfers.POST("/:id/cancel",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.CancelTransfer,
			)

			transfers.POST("/:id/schedule",
				middleware.RequirePermission(middleware.PermissionUpdateTransfers),
				transferHandler.ScheduleTransfer,
			)

			// Initiate transfer
			transfers.POST("/initiate",
				middleware.RequirePermission(middleware.PermissionCreateTransfers),
				transferHandler.InitiateTransfer,
			)
		}

		// Inventory management routes
		inventory := protected.Group("/inventory")
		{
			// List inventory items
			inventory.GET("",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.ListInventoryItems,
			)

			// Get low stock items
			inventory.GET("/low-stock",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetLowStockItems,
			)

			// Get out of stock items
			inventory.GET("/out-of-stock",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetOutOfStockItems,
			)

			// Get expired items
			inventory.GET("/expired",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetExpiredItems,
			)

			// Get expiring soon items
			inventory.GET("/expiring-soon",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetExpiringSoonItems,
			)

			// Get items needing reorder
			inventory.GET("/needing-reorder",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetItemsNeedingReorder,
			)


			// Get inventory statistics
			inventory.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetInventoryStatistics,
			)

			// Get item by SKU
			inventory.GET("/sku/:sku",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetInventoryItemBySKU,
			)

			// Get inventory item by ID
			inventory.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetInventoryItem,
			)

			// Create inventory item
			inventory.POST("",
				middleware.RequirePermission(middleware.PermissionCreateInventory),
				inventoryHandler.CreateInventoryItem,
			)

			// Update inventory item
			inventory.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.UpdateInventoryItem,
			)

			// Delete inventory item
			inventory.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteInventory),
				inventoryHandler.DeleteInventoryItem,
			)

			// Stock operations
			inventory.POST("/:id/add-stock",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.AddStock,
			)

			inventory.POST("/:id/remove-stock",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.RemoveStock,
			)

			inventory.POST("/:id/adjust-stock",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.AdjustStock,
			)

			// Item activation
			inventory.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.ActivateItem,
			)

			inventory.POST("/:id/deactivate",
				middleware.RequirePermission(middleware.PermissionUpdateInventory),
				inventoryHandler.DeactivateItem,
			)

			// Get inventory history
			inventory.GET("/:id/history",
				middleware.RequirePermission(middleware.PermissionViewInventory),
				inventoryHandler.GetInventoryHistory,
			)
		}

		// Stock transaction management routes
		stockTransactions := protected.Group("/stock-transactions")
		{
			// List stock transactions
			stockTransactions.GET("",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.ListStockTransactions,
			)

			// Get transactions by item
			stockTransactions.GET("/item/:item_id",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.GetStockTransactionsByItem,
			)

			// Get transactions by type
			stockTransactions.GET("/type/:type",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.GetStockTransactionsByType,
			)

			// Get transaction statistics
			stockTransactions.GET("/statistics",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.GetStockTransactionStatistics,
			)

			// Get transaction by ID
			stockTransactions.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.GetStockTransaction,
			)

			// Export stock transactions
			stockTransactions.GET("/export",
				middleware.RequirePermission(middleware.PermissionViewStockTransactions),
				stockTransactionHandler.ExportStockTransactions,
			)
		}

		// Audit Log management routes (admin only)
		auditLogs := protected.Group("/audit-logs")
		auditLogs.Use(middleware.RequireAdmin()) // All audit log endpoints require admin role
		{
			// List audit logs with filtering
			auditLogs.GET("", auditLogHandler.ListAuditLogs)

			// Get audit log by ID
			auditLogs.GET("/:id", auditLogHandler.GetAuditLog)

			// Get recent activity (last N logs)
			auditLogs.GET("/recent", auditLogHandler.GetRecentActivity)

			// Get audit statistics
			auditLogs.GET("/statistics", auditLogHandler.GetAuditStatistics)

			// Get logs for a specific user
			auditLogs.GET("/user/:user_id", auditLogHandler.GetUserActivity)

			// Get logs for a specific action type
			auditLogs.GET("/action/:action", auditLogHandler.GetActionLogs)

			// Get logs for a specific entity
			auditLogs.GET("/entity/:entity_type/:entity_id", auditLogHandler.GetEntityHistory)

			// Get logs for date range
			auditLogs.GET("/date-range", auditLogHandler.GetLogsForDateRange)

			// Export audit logs
			auditLogs.GET("/export", auditLogHandler.ExportAuditLogs)

			// Search audit logs
			auditLogs.GET("/search", auditLogHandler.SearchAuditLogs)

			// Delete old logs (super admin only)
			auditLogs.DELETE("/cleanup",
				middleware.RequireSuperAdmin(),
				auditLogHandler.DeleteOldLogs,
			)
		}

		// System Monitoring routes (admin only)
		monitoring := protected.Group("/monitoring")
		monitoring.Use(middleware.RequireAdmin()) // All monitoring endpoints require admin role
		{
			// Get system health
			monitoring.GET("/health", monitoringHandler.GetSystemHealth)

			// Get usage statistics
			monitoring.GET("/statistics", monitoringHandler.GetUsageStatistics)

			// Get performance metrics
			monitoring.GET("/performance", monitoringHandler.GetPerformanceMetrics)

			// Get database statistics
			monitoring.GET("/database", monitoringHandler.GetDatabaseStatistics)

			// Get system configuration
			monitoring.GET("/configuration", monitoringHandler.GetSystemConfiguration)
		}

		// Medical Conditions routes
		conditions := protected.Group("/medical-conditions")
		{
			// List medical conditions
			conditions.GET("",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.ListConditions,
			)

			// Get chronic conditions
			conditions.GET("/chronic",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetChronicConditions,
			)

			// Get condition by ID
			conditions.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetCondition,
			)

			// Create condition
			conditions.POST("",
				middleware.RequirePermission(middleware.PermissionCreateVeterinary),
				medicalHandler.CreateCondition,
			)

			// Update condition
			conditions.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.UpdateCondition,
			)

			// Delete condition
			conditions.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteVeterinary),
				medicalHandler.DeleteCondition,
			)

			// Resolve condition
			conditions.POST("/:id/resolve",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.ResolveCondition,
			)
		}

		// Medications routes
		medications := protected.Group("/medications")
		{
			// List medications
			medications.GET("",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.ListMedications,
			)

			// Get medications due for refill
			medications.GET("/due-for-refill",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetMedicationsDueForRefill,
			)

			// Get medications expiring soon
			medications.GET("/expiring-soon",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetExpiringSoonMedications,
			)

			// Get medication by ID
			medications.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetMedication,
			)

			// Create medication
			medications.POST("",
				middleware.RequirePermission(middleware.PermissionCreateVeterinary),
				medicalHandler.CreateMedication,
			)

			// Update medication
			medications.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.UpdateMedication,
			)

			// Delete medication
			medications.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteVeterinary),
				medicalHandler.DeleteMedication,
			)

			// Record medication administration
			medications.POST("/:id/administer",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.RecordMedicationAdministration,
			)

			// Refill medication
			medications.POST("/:id/refill",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.RefillMedication,
			)
		}

		// Treatment Plans routes
		treatmentPlans := protected.Group("/treatment-plans")
		{
			// List treatment plans
			treatmentPlans.GET("",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.ListTreatmentPlans,
			)

			// Get treatment plan by ID
			treatmentPlans.GET("/:id",
				middleware.RequirePermission(middleware.PermissionViewVeterinary),
				medicalHandler.GetTreatmentPlan,
			)

			// Create treatment plan
			treatmentPlans.POST("",
				middleware.RequirePermission(middleware.PermissionCreateVeterinary),
				medicalHandler.CreateTreatmentPlan,
			)

			// Update treatment plan
			treatmentPlans.PUT("/:id",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.UpdateTreatmentPlan,
			)

			// Delete treatment plan
			treatmentPlans.DELETE("/:id",
				middleware.RequirePermission(middleware.PermissionDeleteVeterinary),
				medicalHandler.DeleteTreatmentPlan,
			)

			// Add progress note
			treatmentPlans.POST("/:id/progress",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.AddProgressNote,
			)

			// Activate treatment plan
			treatmentPlans.POST("/:id/activate",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.ActivateTreatmentPlan,
			)

			// Complete treatment plan
			treatmentPlans.POST("/:id/complete",
				middleware.RequirePermission(middleware.PermissionUpdateVeterinary),
				medicalHandler.CompleteTreatmentPlan,
			)
		}

		// Batch operations routes
		batches := protected.Group("/batches")
		{
			// Create batch operation
			batches.POST("",
				middleware.RequirePermission(middleware.PermissionUpdateCommunications),
				batchHandler.CreateBatch,
			)
		}
	}

	// Public animal routes (no auth required)
	publicAnimals := router.Group("/api/v1/public/animals")
	{
		// Get species list
		publicAnimals.GET("/species", animalHandler.GetSpecies)

		// List available animals (for adoption page)
		publicAnimals.GET("", animalHandler.ListAnimals)
		publicAnimals.GET("/:id", animalHandler.GetAnimal)
	}
}
