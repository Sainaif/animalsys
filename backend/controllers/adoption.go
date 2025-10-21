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

type AdoptionController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterAdoptionRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &AdoptionController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", ctrl.GetAll)
	rg.GET("/:id", ctrl.GetByID)
	rg.POST("", ctrl.Create)
	rg.PUT("/:id", middlewares.RequireRole("admin", "employee"), ctrl.Update)
	rg.DELETE("/:id", middlewares.RequireRole("admin"), ctrl.Delete)
}

func (ctrl *AdoptionController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("adoptions")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch adoptions")
		return
	}
	defer cursor.Close(ctx)

	var adoptions []models.Adoption
	if err = cursor.All(ctx, &adoptions); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode adoptions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adoptions)
}

func (ctrl *AdoptionController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid adoption ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var adoption models.Adoption
	collection := ctrl.db.Collection("adoptions")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&adoption)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Adoption not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adoption)
}

func (ctrl *AdoptionController) Create(c *gin.Context) {
	var adoption models.Adoption
	if err := c.ShouldBindJSON(&adoption); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	adoption.ID = primitive.NewObjectID()
	if adoption.Status == "" {
		adoption.Status = "pending"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("adoptions")
	_, err := collection.InsertOne(ctx, adoption)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create adoption")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, adoption)
}

func (ctrl *AdoptionController) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid adoption ID")
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

	collection := ctrl.db.Collection("adoptions")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update adoption")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Adoption not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Adoption updated successfully"})
}

func (ctrl *AdoptionController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid adoption ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("adoptions")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete adoption")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Adoption not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Adoption deleted successfully"})
}
