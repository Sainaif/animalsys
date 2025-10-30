package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sainaif/animalsys/internal/adapters/auth"
	httpAdapter "github.com/sainaif/animalsys/internal/adapters/http"
	"github.com/sainaif/animalsys/internal/adapters/http/handlers"
	"github.com/sainaif/animalsys/internal/adapters/repository"
	"github.com/sainaif/animalsys/internal/core/usecases"
	"github.com/sainaif/animalsys/internal/infrastructure"
	"github.com/sainaif/animalsys/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/internal/pkg"
)

func main() {
	// Load configuration
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := infrastructure.NewLogger(cfg.Environment, cfg.LogLevel, "")
	logger.Info("Starting AnimalSys ERP System...")

	// Connect to database
	db, err := infrastructure.NewDatabase(cfg, logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Disconnect(context.Background())

	// Health check
	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal(fmt.Sprintf("Database health check failed: %v", err))
	}
	logger.Info("Database connection established")

	// Initialize utilities
	jwtManager := pkg.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	passwordManager := pkg.NewPasswordManager()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.GetClient())
	animalRepo := repository.NewAnimalRepository(db.GetClient())
	adoptionRepo := repository.NewAdoptionRepository(db.GetClient())
	volunteerRepo := repository.NewVolunteerRepository(db.GetClient())
	volunteerHourRepo := repository.NewVolunteerHourRepository(db.GetClient())
	scheduleRepo := repository.NewScheduleRepository(db.GetClient())
	documentRepo := repository.NewDocumentRepository(db.GetClient())
	financeRepo := repository.NewFinanceRepository(db.GetClient())
	donorRepo := repository.NewDonorRepository(db.GetClient())
	donationRepo := repository.NewDonationRepository(db.GetClient())
	inventoryRepo := repository.NewInventoryRepository(db.GetClient())
	stockMovementRepo := repository.NewStockMovementRepository(db.GetClient())
	veterinaryRepo := repository.NewVeterinaryRepository(db.GetClient())
	vaccinationRepo := repository.NewVaccinationRepository(db.GetClient())
	campaignRepo := repository.NewCampaignRepository(db.GetClient())
	partnerRepo := repository.NewPartnerRepository(db.GetClient())
	communicationRepo := repository.NewCommunicationRepository(db.GetClient())
	templateRepo := repository.NewCommunicationTemplateRepository(db.GetClient())
	auditRepo := repository.NewAuditLogRepository(db.GetClient())

	logger.Info("Repositories initialized")

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo, auditRepo, jwtManager, passwordManager)
	userUseCase := usecases.NewUserUseCase(userRepo, auditRepo)
	animalUseCase := usecases.NewAnimalUseCase(animalRepo, auditRepo)
	adoptionUseCase := usecases.NewAdoptionUseCase(adoptionRepo, animalRepo, auditRepo)
	volunteerUseCase := usecases.NewVolunteerUseCase(volunteerRepo, volunteerHourRepo, auditRepo)
	scheduleUseCase := usecases.NewScheduleUseCase(scheduleRepo, userRepo, auditRepo)
	documentUseCase := usecases.NewDocumentUseCase(documentRepo, auditRepo)
	financeUseCase := usecases.NewFinanceUseCase(financeRepo, auditRepo)
	donorUseCase := usecases.NewDonorUseCase(donorRepo, donationRepo, auditRepo)
	inventoryUseCase := usecases.NewInventoryUseCase(inventoryRepo, stockMovementRepo, auditRepo)
	veterinaryUseCase := usecases.NewVeterinaryUseCase(veterinaryRepo, vaccinationRepo, animalRepo, auditRepo)
	campaignUseCase := usecases.NewCampaignUseCase(campaignRepo, auditRepo)
	partnerUseCase := usecases.NewPartnerUseCase(partnerRepo, auditRepo)
	communicationUseCase := usecases.NewCommunicationUseCase(communicationRepo, templateRepo, auditRepo)

	logger.Info("Use cases initialized")

	// Initialize auth middleware and RBAC
	authMiddleware := auth.NewAuthMiddleware(jwtManager)
	rbac := auth.NewRBAC()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)
	animalHandler := handlers.NewAnimalHandler(animalUseCase)
	adoptionHandler := handlers.NewAdoptionHandler(adoptionUseCase)
	volunteerHandler := handlers.NewVolunteerHandler(volunteerUseCase)
	scheduleHandler := handlers.NewScheduleHandler(scheduleUseCase)
	documentHandler := handlers.NewDocumentHandler(documentUseCase)
	financeHandler := handlers.NewFinanceHandler(financeUseCase)
	donorHandler := handlers.NewDonorHandler(donorUseCase)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUseCase)
	veterinaryHandler := handlers.NewVeterinaryHandler(veterinaryUseCase)
	campaignHandler := handlers.NewCampaignHandler(campaignUseCase)
	partnerHandler := handlers.NewPartnerHandler(partnerUseCase)
	communicationHandler := handlers.NewCommunicationHandler(communicationUseCase)

	logger.Info("Handlers initialized")

	// Initialize middleware
	corsMiddleware := middleware.CORS(cfg.CORSOrigins)
	securityMiddleware := middleware.SecurityHeaders()
	loggerMiddleware := middleware.Logger(logger)
	rateLimitMiddleware := middleware.RateLimit(100, 15*time.Minute)

	// Initialize router
	router := httpAdapter.NewRouter(
		authHandler,
		userHandler,
		animalHandler,
		adoptionHandler,
		volunteerHandler,
		scheduleHandler,
		documentHandler,
		financeHandler,
		donorHandler,
		inventoryHandler,
		veterinaryHandler,
		campaignHandler,
		partnerHandler,
		communicationHandler,
		authMiddleware,
		rbac,
	)

	router.SetupRoutes(corsMiddleware, securityMiddleware, loggerMiddleware, rateLimitMiddleware)

	logger.Info(fmt.Sprintf("Server configured to run on port %s", cfg.ServerPort))

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.ServerPort)
		if err := router.Run(addr); err != nil {
			logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	logger.Info("Server started successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Disconnect(ctx); err != nil {
		logger.Error(fmt.Sprintf("Error disconnecting from database: %v", err))
	}

	logger.Info("Server shutdown complete")
}
