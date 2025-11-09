package auth

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/sainaif/animalsys/backend/pkg/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo        repositories.UserRepository
	auditLogRepo    repositories.AuditLogRepository
	jwtService      *security.JWTService
	passwordService *security.PasswordService
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	auditLogRepo repositories.AuditLogRepository,
	jwtService *security.JWTService,
	passwordService *security.PasswordService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		auditLogRepo:    auditLogRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	User         *entities.User `json:"user"`
}

// Login authenticates a user and returns tokens
func (uc *AuthUseCase) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !uc.passwordService.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.NewForbidden("user account is not active")
	}

	// Generate access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	if err := uc.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(user.ID, entities.ActionLogin, "user", ipAddress, userAgent)
	_ = uc.auditLogRepo.Create(ctx, auditLog) // Don't fail login if audit log fails

	// Remove sensitive data before returning
	user.PasswordHash = ""
	user.RefreshToken = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken refreshes the access token using a refresh token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// Validate refresh token
	userID, err := uc.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Get user from database
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Verify that the refresh token matches the one in database
	if user.RefreshToken != req.RefreshToken {
		return nil, errors.ErrInvalidToken
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.NewForbidden("user account is not active")
	}

	// Generate new access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Optionally rotate refresh token (more secure)
	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update refresh token in database
	if err := uc.userRepo.UpdateRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout logs out a user by invalidating their refresh token
func (uc *AuthUseCase) Logout(ctx context.Context, userID primitive.ObjectID, ipAddress, userAgent string) error {
	// Clear refresh token
	if err := uc.userRepo.UpdateRefreshToken(ctx, userID, ""); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionLogout, "user", ipAddress, userAgent)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string             `json:"email" validate:"required,email"`
	Password  string             `json:"password" validate:"required,min=8"`
	FirstName string             `json:"first_name" validate:"required"`
	LastName  string             `json:"last_name" validate:"required"`
	Role      entities.UserRole  `json:"role" validate:"required"`
	Language  string             `json:"language" validate:"required,oneof=en pl"`
	Theme     string             `json:"theme" validate:"required,oneof=light dark"`
}

// Register creates a new user (admin only)
func (uc *AuthUseCase) Register(ctx context.Context, req *RegisterRequest, creatorID primitive.ObjectID) (*entities.User, error) {
	// Validate password strength
	if err := uc.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Status:       entities.StatusActive,
		Language:     req.Language,
		Theme:        req.Theme,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "user", "", "").
		WithEntityID(user.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	// Remove sensitive data
	user.PasswordHash = ""

	return user, nil
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePassword changes a user's password
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID primitive.ObjectID, req *ChangePasswordRequest) error {
	// Get user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !uc.passwordService.VerifyPassword(req.OldPassword, user.PasswordHash) {
		return errors.NewBadRequest("incorrect current password")
	}

	// Validate new password strength
	if err := uc.passwordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := uc.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Clear refresh token to force re-login on all devices
	_ = uc.userRepo.UpdateRefreshToken(ctx, userID, "")

	// Create audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "user", "", "").
		WithEntityID(userID).
		WithChanges(map[string]interface{}{"password": "changed"})
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetCurrentUser returns the current user details
func (uc *AuthUseCase) GetCurrentUser(ctx context.Context, userID primitive.ObjectID) (*entities.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.RefreshToken = ""

	return user, nil
}
