package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/auth"
	"github.com/sainaif/animalsys/backend/pkg/errors"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUseCase *auth.AuthUseCase
	validate    *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validate:    validator.New(),
	}
}

// Login handles user login
// @Summary Login
// @Description Authenticate user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authUseCase.Login(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} auth.RefreshTokenResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
// @Summary Logout
// @Description Logout user and invalidate refresh token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} errors.AppError
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.authUseCase.Logout(c.Request.Context(), *userID, ipAddress, userAgent); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Register handles user registration (admin only)
// @Summary Register User
// @Description Create a new user account (admin only)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "User details"
// @Success 201 {object} entities.User
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req auth.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), &req, *userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetMe returns current user details
// @Summary Get Current User
// @Description Get details of currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entities.User
// @Failure 401 {object} errors.AppError
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authUseCase.GetCurrentUser(c.Request.Context(), *userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword handles password change
// @Summary Change Password
// @Description Change user's password
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body auth.ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req auth.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUseCase.ChangePassword(c.Request.Context(), *userID, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
