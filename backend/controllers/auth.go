package controllers

import (
	"animalsys/config"
	"animalsys/models"
	"animalsys/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	db  *mongo.Database
	cfg config.Config
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

func RegisterAuthRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &AuthController{db: db, cfg: cfg}
	rg.POST("/login", ctrl.Login)
	rg.POST("/register", ctrl.Register)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	collection := ctrl.db.Collection("users")

	if ctrl.cfg.LDAPEnabled {
		ldapConfig := utils.LDAPConfig{
			Server:       ctrl.cfg.LDAPServer,
			BaseDN:       ctrl.cfg.LDAPBaseDN,
			BindDN:       ctrl.cfg.LDAPBindDN,
			BindPassword: ctrl.cfg.LDAPBindPassword,
		}

		ldapID, err := utils.AuthenticateLDAP(req.Username, req.Password, ldapConfig)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		err = collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			user = models.User{
				ID:       primitive.NewObjectID(),
				Username: req.Username,
				Email:    req.Username + "@ldap",
				Role:     "user",
				LDAPID:   &ldapID,
			}
			_, err = collection.InsertOne(ctx, user)
			if err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
				return
			}
		} else if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Database error")
			return
		}
	} else {
		err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, user.Role, ctrl.cfg.JWTSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("users")

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}
