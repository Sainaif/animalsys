package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sainaif/animalsys/backend/internal/delivery/http/handlers"
	"github.com/sainaif/animalsys/backend/internal/delivery/http/routes"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/config"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/logger"
	adoptionUC "github.com/sainaif/animalsys/backend/internal/usecase/adoption"
	animalUC "github.com/sainaif/animalsys/backend/internal/usecase/animal"
	auditlogUC "github.com/sainaif/animalsys/backend/internal/usecase/auditlog"
	authUC "github.com/sainaif/animalsys/backend/internal/usecase/auth"
	campaignUC "github.com/sainaif/animalsys/backend/internal/usecase/campaign"
	communicationUC "github.com/sainaif/animalsys/backend/internal/usecase/communication"
	contactUC "github.com/sainaif/animalsys/backend/internal/usecase/contact"
	dashboardUC "github.com/sainaif/animalsys/backend/internal/usecase/dashboard"
	documentUC "github.com/sainaif/animalsys/backend/internal/usecase/document"
	donationUC "github.com/sainaif/animalsys/backend/internal/usecase/donation"
	donorUC "github.com/sainaif/animalsys/backend/internal/usecase/donor"
	eventUC "github.com/sainaif/animalsys/backend/internal/usecase/event"
	inventoryUC "github.com/sainaif/animalsys/backend/internal/usecase/inventory"
	medicalUC "github.com/sainaif/animalsys/backend/internal/usecase/medical"
	monitoringUC "github.com/sainaif/animalsys/backend/internal/usecase/monitoring"
	notificationUC "github.com/sainaif/animalsys/backend/internal/usecase/notification"
	partnerUC "github.com/sainaif/animalsys/backend/internal/usecase/partner"
	reportUC "github.com/sainaif/animalsys/backend/internal/usecase/report"
	settingsUC "github.com/sainaif/animalsys/backend/internal/usecase/settings"
	stockUC "github.com/sainaif/animalsys/backend/internal/usecase/stock"
	taskUC "github.com/sainaif/animalsys/backend/internal/usecase/task"
	transferUC "github.com/sainaif/animalsys/backend/internal/usecase/transfer"
	userUC "github.com/sainaif/animalsys/backend/internal/usecase/user"
	veterinaryUC "github.com/sainaif/animalsys/backend/internal/usecase/veterinary"
	volunteerUC "github.com/sainaif/animalsys/backend/internal/usecase/volunteer"
	"github.com/sainaif/animalsys/backend/pkg/security"
	"github.com/sainaif/animalsys/backend/pkg/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	logger.Init(cfg.Log.Level, cfg.Environment)

	log.Info().
		Str("environment", cfg.Environment).
		Str("version", cfg.Version).
		Msg("Starting Animal Foundation CRM")

	// Initialize database
	ctx := context.Background()
	db, err := mongodb.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			log.Error().Err(err).Msg("Error disconnecting from database")
		}
	}()

	log.Info().Msg("Successfully connected to MongoDB")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	animalRepo := repositories.NewAnimalRepository(db)
	veterinaryVisitRepo := repositories.NewVeterinaryVisitRepository(db)
	vaccinationRepo := repositories.NewVaccinationRepository(db)
	adoptionApplicationRepo := repositories.NewAdoptionApplicationRepository(db)
	adoptionRepo := repositories.NewAdoptionRepository(db)
	donorRepo := repositories.NewDonorRepository(db)
	donationRepo := repositories.NewDonationRepository(db)
	campaignRepo := repositories.NewCampaignRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	volunteerRepo := repositories.NewVolunteerRepository(db)
	contactRepo := repositories.NewContactRepository(db)
	eventAttendanceRepo := repositories.NewEventAttendanceRepository(db)
	volunteerAssignmentRepo := repositories.NewVolunteerAssignmentRepository(db)
	communicationTemplateRepo := repositories.NewCommunicationTemplateRepository(db)
	communicationRepo := repositories.NewCommunicationRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	reportExecutionRepo := repositories.NewReportExecutionRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	documentRepo := repositories.NewDocumentRepository(db)
	partnerRepo := repositories.NewPartnerRepository(db)
	transferRepo := repositories.NewTransferRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	stockTransactionRepo := repositories.NewStockTransactionRepository(db)
	medicalConditionRepo := repositories.NewMedicalConditionRepository(db)
	medicationRepo := repositories.NewMedicationRepository(db)
	treatmentPlanRepo := repositories.NewTreatmentPlanRepository(db)

	// Ensure database indexes
	if err := userRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create user indexes")
	}
	if err := auditLogRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create audit log indexes")
	}
	if err := animalRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create animal indexes")
	}
	if err := veterinaryVisitRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create veterinary visit indexes")
	}
	if err := vaccinationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create vaccination indexes")
	}
	if err := adoptionApplicationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create adoption application indexes")
	}
	if err := adoptionRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create adoption indexes")
	}
	if err := donorRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create donor indexes")
	}
	if err := donationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create donation indexes")
	}
	if err := campaignRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create campaign indexes")
	}
	if err := eventRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create event indexes")
	}
	if err := contactRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create contact indexes")
	}
	if err := volunteerRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create volunteer indexes")
	}
	if err := eventAttendanceRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create event attendance indexes")
	}
	if err := volunteerAssignmentRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create volunteer assignment indexes")
	}
	if err := communicationTemplateRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create communication template indexes")
	}
	if err := communicationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create communication indexes")
	}
	if err := notificationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create notification indexes")
	}
	if err := reportRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create report indexes")
	}
	if err := reportExecutionRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create report execution indexes")
	}
	if err := settingsRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create settings indexes")
	}
	if err := taskRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create task indexes")
	}
	if err := documentRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create document indexes")
	}
	if err := partnerRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create partner indexes")
	}
	if err := transferRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create transfer indexes")
	}
	if err := inventoryRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create inventory indexes")
	}
	if err := stockTransactionRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create stock transaction indexes")
	}
	if err := medicalConditionRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create medical condition indexes")
	}
	if err := medicationRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create medication indexes")
	}
	if err := treatmentPlanRepo.EnsureIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to create treatment plan indexes")
	}

	// Initialize security services
	jwtService := security.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	passwordService := security.NewPasswordService()

	// Initialize storage service
	storageService := storage.NewStorageService(
		cfg.Storage.LocalPath,
		cfg.Storage.BaseURL,
		cfg.Storage.MaxFileSize,
	)

	// Initialize use cases
	authUseCase := authUC.NewAuthUseCase(
		userRepo,
		auditLogRepo,
		jwtService,
		passwordService,
	)
	userUseCase := userUC.NewUserUseCase(
		userRepo,
		auditLogRepo,
		passwordService,
	)
	animalUseCase := animalUC.NewAnimalUseCase(
		animalRepo,
		auditLogRepo,
		storageService,
	)
	veterinaryUseCase := veterinaryUC.NewVeterinaryUseCase(
		veterinaryVisitRepo,
		vaccinationRepo,
		animalRepo,
		auditLogRepo,
	)
	adoptionUseCase := adoptionUC.NewAdoptionUseCase(
		adoptionApplicationRepo,
		adoptionRepo,
		animalRepo,
		auditLogRepo,
	)
	donorUseCase := donorUC.NewDonorUseCase(
		donorRepo,
		auditLogRepo,
	)
	donationUseCase := donationUC.NewDonationUseCase(
		donationRepo,
		donorRepo,
		campaignRepo,
		auditLogRepo,
	)
	campaignUseCase := campaignUC.NewCampaignUseCase(
		campaignRepo,
		auditLogRepo,
	)
	eventUseCase := eventUC.NewEventUseCase(
		eventRepo,
		eventAttendanceRepo,
		volunteerRepo,
		auditLogRepo,
	)
	volunteerUseCase := volunteerUC.NewVolunteerUseCase(
		volunteerRepo,
		volunteerAssignmentRepo,
		auditLogRepo,
	)
	contactUseCase := contactUC.NewUseCase(contactRepo)
	communicationUseCase := communicationUC.NewCommunicationUseCase(
		communicationRepo,
		communicationTemplateRepo,
		auditLogRepo,
	)
	notificationUseCase := notificationUC.NewNotificationUseCase(
		notificationRepo,
		auditLogRepo,
	)
	reportUseCase := reportUC.NewReportUseCase(
		reportRepo,
		reportExecutionRepo,
		auditLogRepo,
	)
	dashboardUseCase := dashboardUC.NewDashboardUseCase(
		animalRepo,
		adoptionRepo,
		donationRepo,
		volunteerRepo,
	)
	settingsUseCase := settingsUC.NewSettingsUseCase(
		settingsRepo,
		auditLogRepo,
	)
	taskUseCase := taskUC.NewTaskUseCase(
		taskRepo,
		auditLogRepo,
	)
	documentUseCase := documentUC.NewDocumentUseCase(
		documentRepo,
		auditLogRepo,
	)
	partnerUseCase := partnerUC.NewPartnerUseCase(
		partnerRepo,
		auditLogRepo,
	)
	transferUseCase := transferUC.NewTransferUseCase(
		transferRepo,
		animalRepo,
		partnerRepo,
		auditLogRepo,
	)
	inventoryUseCase := inventoryUC.NewInventoryUseCase(
		inventoryRepo,
		stockTransactionRepo,
		auditLogRepo,
	)
	stockTransactionUseCase := stockUC.NewStockTransactionUseCase(
		stockTransactionRepo,
		inventoryRepo,
	)
	auditLogUseCase := auditlogUC.NewAuditLogUseCase(
		auditLogRepo,
		userRepo,
	)
	monitoringUseCase := monitoringUC.NewMonitoringUseCase(
		db,
		userRepo,
		animalRepo,
		adoptionRepo,
		donationRepo,
		documentRepo,
		auditLogRepo,
	)
	medicalUseCase := medicalUC.NewMedicalUseCase(
		medicalConditionRepo,
		medicationRepo,
		treatmentPlanRepo,
		animalRepo,
		auditLogRepo,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)
	animalHandler := handlers.NewAnimalHandler(animalUseCase)
	veterinaryHandler := handlers.NewVeterinaryHandler(veterinaryUseCase)
	adoptionHandler := handlers.NewAdoptionHandler(adoptionUseCase)
	donorHandler := handlers.NewDonorHandler(donorUseCase)
	donationHandler := handlers.NewDonationHandler(donationUseCase)
	campaignHandler := handlers.NewCampaignHandler(campaignUseCase)
	eventHandler := handlers.NewEventHandler(eventUseCase)
	volunteerHandler := handlers.NewVolunteerHandler(volunteerUseCase)
	contactHandler := handlers.NewContactHandler(contactUseCase)
	communicationHandler := handlers.NewCommunicationHandler(communicationUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	reportHandler := handlers.NewReportHandler(reportUseCase)
	dashboardHandler := handlers.NewDashboardHandler(
		dashboardUseCase,
		animalUseCase,
		taskUseCase,
	)
	settingsHandler := handlers.NewSettingsHandler(settingsUseCase)
	taskHandler := handlers.NewTaskHandler(taskUseCase)
	documentHandler := handlers.NewDocumentHandler(documentUseCase)
	partnerHandler := handlers.NewPartnerHandler(partnerUseCase)
	transferHandler := handlers.NewTransferHandler(transferUseCase)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUseCase)
	stockTransactionHandler := handlers.NewStockTransactionHandler(stockTransactionUseCase)
	auditLogHandler := handlers.NewAuditLogHandler(auditLogUseCase)
	monitoringHandler := handlers.NewMonitoringHandler(monitoringUseCase)
	medicalHandler := handlers.NewMedicalHandler(medicalUseCase)
	batchHandler := handlers.NewBatchHandler()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))
	router.Use(loggerMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheckHandler(db))

	// Setup API routes
	routes.SetupRoutes(router, authHandler, userHandler, animalHandler, veterinaryHandler, adoptionHandler, donorHandler, donationHandler, campaignHandler, eventHandler, volunteerHandler, contactHandler, communicationHandler, notificationHandler, reportHandler, dashboardHandler, settingsHandler, taskHandler, documentHandler, partnerHandler, transferHandler, inventoryHandler, stockTransactionHandler, auditLogHandler, monitoringHandler, medicalHandler, batchHandler, jwtService, userRepo)

	// Create server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Info().
			Str("port", cfg.Server.Port).
			Msg("Starting HTTP server")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server stopped gracefully")
}

// corsMiddleware handles CORS
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := parseAllowedOrigins(cfg.CORS.AllowedOrigins)
	allowAll := len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowAll {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin != "" && containsOrigin(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return nil
	}

	items := strings.Split(origins, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func containsOrigin(origins []string, origin string) bool {
	for _, allowed := range origins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// loggerMiddleware logs HTTP requests using zerolog
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		var logEvent *zerolog.Event
		if statusCode >= 500 {
			logEvent = log.Error()
		} else if statusCode >= 400 {
			logEvent = log.Warn()
		} else {
			logEvent = log.Info()
		}

		logEvent.
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("error", errorMessage).
			Msg("HTTP request")
	}
}

// healthCheckHandler returns server health status
func healthCheckHandler(db *mongodb.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			log.Error().Err(err).Msg("Database health check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
			"time":     time.Now().UTC(),
		})
	}
}
