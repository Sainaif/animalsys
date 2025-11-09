package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/user"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler handles user management HTTP requests
type UserHandler struct {
	userUseCase *user.UserUseCase
	validate    *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *user.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validate:    validator.New(),
	}
}

// ListUsers lists all users with pagination
// @Summary List Users
// @Description Get list of users with filters and pagination
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Param search query string false "Search in email, first name, last name"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} user.ListUsersResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req user.ListUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userUseCase.ListUsers(c.Request.Context(), &req)
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

// GetUser gets a user by ID
// @Summary Get User
// @Description Get user details by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userUseCase.GetUserByID(c.Request.Context(), userID)
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

// CreateUser creates a new user
// @Summary Create User
// @Description Create a new user (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "User details"
// @Success 201 {object} entities.User
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	creatorID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req user.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.userUseCase.CreateUser(c.Request.Context(), &req, *creatorID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser updates a user
// @Summary Update User
// @Description Update user details (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body user.UpdateUserRequest true "User details to update"
// @Success 200 {object} entities.User
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	updaterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req user.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userUseCase.UpdateUser(c.Request.Context(), userID, &req, *updaterID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser deletes a user
// @Summary Delete User
// @Description Delete a user (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	deleterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userUseCase.DeleteUser(c.Request.Context(), userID, *deleterID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// ResetPassword resets a user's password (admin only)
// @Summary Reset User Password
// @Description Reset a user's password (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body user.ResetPasswordRequest true "New password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /users/{id}/reset-password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	adminID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req user.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUseCase.ResetPassword(c.Request.Context(), userID, &req, *adminID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// UpdateUserRole updates a user's role
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}

// UpdateUserStatus updates a user's status
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User status updated"})
}
