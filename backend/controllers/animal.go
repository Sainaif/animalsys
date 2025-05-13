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

func RegisterAnimalRoutes(rg *gin.RouterGroup, db *mongo.Database) {
	animals := db.Collection("animals")

	rg.GET("/", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cur, err := animals.Find(ctx, bson.M{})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.Animal
		if err := cur.All(ctx, &result); err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, result)
	})

	rg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var animal models.Animal
		err = animals.FindOne(ctx, bson.M{"_id": objID}).Decode(&animal)
		if err != nil {
			utils.Error(c, http.StatusNotFound, "Animal not found")
			return
		}
		utils.Success(c, animal)
	})

	rg.POST("/", middlewares.AuthMiddleware, middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		var animal models.Animal
		if err := c.ShouldBindJSON(&animal); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		animal.Status = "available"
		res, err := animals.InsertOne(ctx, animal)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		animal.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, animal)
	})

	rg.PUT("/:id", middlewares.AuthMiddleware, middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var animal models.Animal
		if err := c.ShouldBindJSON(&animal); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		update := bson.M{
			"name":           animal.Name,
			"species":        animal.Species,
			"breed":          animal.Breed,
			"age":            animal.Age,
			"health_history": animal.HealthHistory,
			"status":         animal.Status,
		}
		_, err = animals.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "updated")
	})

	rg.DELETE("/:id", middlewares.AuthMiddleware, middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = animals.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "deleted")
	})
}
