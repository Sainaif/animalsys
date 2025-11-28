package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/pkg/security"
)

// TestHandler handles test-only HTTP requests
type TestHandler struct {
	jwtService *security.JWTService
	userRepo   repositories.UserRepository
}

// NewTestHandler creates a new test handler
func NewTestHandler(jwtService *security.JWTService, userRepo repositories.UserRepository) *TestHandler {
	return &TestHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// GenerateTestToken generates a test token with a custom expiration
// @Summary Generate Test Token
// @Description Generate a test token with a custom expiration
// @Tags test
// @Security BearerAuth
// @Produce json
// @Param expires_in query int true "Expiration time in seconds"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /test/generate-token [post]
func (h *TestHandler) GenerateTestToken(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ExpiresIn int `json:"expires_in" form:"expires_in" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	user, err := h.userRepo.FindByID(c.Request.Context(), *userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Email, string(user.Role), time.Duration(req.ExpiresIn)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
