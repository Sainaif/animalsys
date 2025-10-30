package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sainaif/animalsys/internal/adapters/repository"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/infrastructure"
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
	logger.Info("Starting database seed...")

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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.GetClient())
	auditRepo := repository.NewAuditLogRepository(db.GetClient())

	// Initialize password manager
	passwordManager := pkg.NewPasswordManager()

	ctx := context.Background()

	// Check if super admin already exists
	existingUsers, _, err := userRepo.List(ctx, &entities.UserFilter{
		Role:   entities.RoleSuperAdmin,
		Limit:  1,
		Offset: 0,
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to check existing users: %v", err))
	}

	if len(existingUsers) > 0 {
		logger.Info("Super admin already exists. Skipping seed.")
		return
	}

	// Create super admin
	hashedPassword, err := passwordManager.HashPassword("Admin123!")
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to hash password: %v", err))
	}

	superAdmin := entities.NewUser(
		"admin@animalsys.local",
		"admin",
		hashedPassword,
		entities.RoleSuperAdmin,
	)
	superAdmin.FirstName = "Super"
	superAdmin.LastName = "Admin"
	superAdmin.Active = true
	superAdmin.CreatedBy = "system"

	if err := userRepo.Create(ctx, superAdmin); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create super admin: %v", err))
	}

	// Create audit log
	auditLog := entities.NewAuditLog(
		"system",
		"",
		"",
		entities.ActionCreate,
		"user",
		superAdmin.ID.Hex(),
		"Initial super admin created",
	)
	if err := auditRepo.Create(ctx, auditLog); err != nil {
		logger.Error(fmt.Sprintf("Failed to create audit log: %v", err))
	}

	logger.Info("========================================")
	logger.Info("Super Admin created successfully!")
	logger.Info("========================================")
	logger.Info(fmt.Sprintf("Email: %s", superAdmin.Email))
	logger.Info(fmt.Sprintf("Username: %s", superAdmin.Username))
	logger.Info("Password: Admin123!")
	logger.Info("========================================")
	logger.Info("IMPORTANT: Change the password after first login!")
	logger.Info("========================================")
}
