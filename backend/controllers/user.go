package controllers

import (
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

func RegisterUserRoutes(rg *gin.RouterGroup, db *mongo.Database) {
	users := db.Collection("users")

	rg.Use(middlewares.AuthMiddleware)

	rg.GET("/", middlewares.RBACMiddleware("admin"), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cur, err := users.Find(ctx, bson.M{})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.User
		if err := cur.All(ctx, &result); err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, result)
	})

	rg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString("user_id")
		role := c.GetString("role")
		if role != "admin" && id != userID {
			utils.Error(c, http.StatusForbidden, "Forbidden")
			return
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user models.User
		err = users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			utils.Error(c, http.StatusNotFound, "User not found")
			return
		}
		utils.Success(c, user)
	})

	rg.POST("/", middlewares.RBACMiddleware("admin"), func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user.PasswordHash, _ = utils.HashPassword(user.PasswordHash)
		res, err := users.InsertOne(ctx, user)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		user.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, user)
	})

	rg.PUT("/:id", middlewares.RBACMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		update := bson.M{
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		}
		if user.PasswordHash != "" {
			hash, _ := utils.HashPassword(user.PasswordHash)
			update["password_hash"] = hash
		}
		_, err = users.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "updated")
	})

	rg.DELETE("/:id", middlewares.RBACMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = users.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "deleted")
	})
}
