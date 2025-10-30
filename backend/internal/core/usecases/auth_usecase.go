package usecases

import (
	"context"
	"errors"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
	"github.com/sainaif/animalsys/internal/pkg"
)

type AuthUseCase struct {
	userRepo        interfaces.UserRepository
	auditRepo       interfaces.AuditLogRepository
	jwtManager      *pkg.JWTManager
	passwordManager *pkg.PasswordManager
}

func NewAuthUseCase(
	userRepo interfaces.UserRepository,
	auditRepo interfaces.AuditLogRepository,
	jwtManager *pkg.JWTManager,
	passwordManager *pkg.PasswordManager,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		auditRepo:       auditRepo,
		jwtManager:      jwtManager,
		passwordManager: passwordManager,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *entities.UserCreateRequest, createdBy string) (*entities.User, error) {
	// Validate password strength
	if err := uc.passwordManager.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, _ = uc.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := uc.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entities.NewUser(
		req.Username,
		req.Email,
		hashedPassword,
		req.Role,
		req.FirstName,
		req.LastName,
	)

	user.Phone = req.Phone
	user.Address = req.Address
	user.CreatedBy = createdBy

	// Save to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Audit log
	auditLog := entities.NewAuditLog(
		user.ID.Hex(),
		user.Email,
		user.Role,
		entities.ActionCreate,
		"user",
		user.ID.Hex(),
		"User registered",
	)
	uc.auditRepo.Create(ctx, auditLog)

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	// Find user by username or email
	var user *entities.User
	var err error

	user, err = uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		user, err = uc.userRepo.GetByEmail(ctx, req.Username)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user account is disabled")
	}

	// Verify password
	if err := uc.passwordManager.ComparePassword(user.Password, req.Password); err != nil {
		// Audit failed login
		auditLog := entities.NewAuditLog(
			user.ID.Hex(),
			user.Email,
			user.Role,
			entities.ActionLogin,
			"auth",
			"",
			"Failed login attempt",
		)
		auditLog.SetError("Invalid password")
		uc.auditRepo.Create(ctx, auditLog)

		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Update last login
	uc.userRepo.UpdateLastLogin(ctx, user.ID.Hex())

	// Audit successful login
	auditLog := entities.NewAuditLog(
		user.ID.Hex(),
		user.Email,
		user.Role,
		entities.ActionLogin,
		"auth",
		"",
		"Successful login",
	)
	uc.auditRepo.Create(ctx, auditLog)

	return &entities.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, req *entities.RefreshTokenRequest) (string, error) {
	// Validate refresh token
	claims, err := uc.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Get user to verify still active
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.Active {
		return "", errors.New("user account is disabled")
	}

	// Generate new access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID string, req *entities.ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := uc.passwordManager.ComparePassword(user.Password, req.OldPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if err := uc.passwordManager.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := uc.passwordManager.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := uc.userRepo.ChangePassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(
		userID,
		user.Email,
		user.Role,
		entities.ActionUpdate,
		"user",
		userID,
		"Password changed",
	)
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}
