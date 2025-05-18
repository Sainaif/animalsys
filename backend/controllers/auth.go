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

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LDAPLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuthRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	rg.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		users := db.Collection("users")
		// Check if username or email exists
		count, _ := users.CountDocuments(ctx, bson.M{"$or": []bson.M{{"username": req.Username}, {"email": req.Email}}})
		if count > 0 {
			utils.Error(c, http.StatusConflict, "Username or email already exists")
			return
		}
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Password hash error")
			return
		}
		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hash,
			Role:         "user",
		}
		res, err := users.InsertOne(ctx, user)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		user.ID = res.InsertedID.(primitive.ObjectID)
		token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, user.Role, cfg.JWTSecret)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "JWT error")
			return
		}
		utils.Success(c, gin.H{"token": token, "role": user.Role})
	})

	rg.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		users := db.Collection("users")
		var user models.User
		err := users.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		ok, err := utils.VerifyPassword(req.Password, user.PasswordHash)
		if err != nil || !ok {
			utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, user.Role, cfg.JWTSecret)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "JWT error")
			return
		}
		utils.Success(c, gin.H{"token": token, "role": user.Role})
	})

	rg.POST("/ldap-login", func(c *gin.Context) {
		var req LDAPLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ldapID, err := utils.LDAPAuthenticate(req.Username, req.Password, cfg)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "LDAP auth failed")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		users := db.Collection("users")
		var user models.User
		err = users.FindOne(ctx, bson.M{"ldap_id": ldapID}).Decode(&user)
		if err != nil {
			// Create user if not exists
			user = models.User{
				Username: req.Username,
				Email:    req.Username + "@ldap.local",
				Role:     "user",
				LDAPID:   &ldapID,
			}
			res, err := users.InsertOne(ctx, user)
			if err != nil {
				utils.Error(c, http.StatusInternalServerError, "DB error")
				return
			}
			user.ID = res.InsertedID.(primitive.ObjectID)
		}
		token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, user.Role, cfg.JWTSecret)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "JWT error")
			return
		}
		utils.Success(c, gin.H{"token": token, "role": user.Role})
	})
}
