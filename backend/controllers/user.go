package controllers

import (
	"animalsys/config"
	"animalsys/middlewares"
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

type UserController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterUserRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &UserController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", ctrl.GetAll)
	rg.GET("/:id", ctrl.GetByID)
	rg.PUT("/:id", ctrl.Update)
	rg.DELETE("/:id", middlewares.RequireRole("admin"), ctrl.Delete)
}

func (ctrl *UserController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode users")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, users)
}

func (ctrl *UserController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	collection := ctrl.db.Collection("users")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	delete(updates, "_id")
	delete(updates, "password_hash")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("users")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("users")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
