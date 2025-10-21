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

type AnimalController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterAnimalRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &AnimalController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", ctrl.GetAll)
	rg.GET("/:id", ctrl.GetByID)
	rg.POST("", middlewares.RequireRole("admin", "employee"), ctrl.Create)
	rg.PUT("/:id", middlewares.RequireRole("admin", "employee"), ctrl.Update)
	rg.DELETE("/:id", middlewares.RequireRole("admin"), ctrl.Delete)
}

func (ctrl *AnimalController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("animals")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch animals")
		return
	}
	defer cursor.Close(ctx)

	var animals []models.Animal
	if err = cursor.All(ctx, &animals); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode animals")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, animals)
}

func (ctrl *AnimalController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid animal ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var animal models.Animal
	collection := ctrl.db.Collection("animals")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&animal)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Animal not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, animal)
}

func (ctrl *AnimalController) Create(c *gin.Context) {
	var animal models.Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	animal.ID = primitive.NewObjectID()
	if animal.HealthHistory == nil {
		animal.HealthHistory = []string{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("animals")
	_, err := collection.InsertOne(ctx, animal)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create animal")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, animal)
}

func (ctrl *AnimalController) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid animal ID")
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	delete(updates, "_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("animals")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update animal")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Animal not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Animal updated successfully"})
}

func (ctrl *AnimalController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid animal ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("animals")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete animal")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Animal not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
