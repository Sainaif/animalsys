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

type FinanceController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterFinanceRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &FinanceController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", middlewares.RequireRole("admin", "employee"), ctrl.GetAll)
	rg.GET("/:id", middlewares.RequireRole("admin", "employee"), ctrl.GetByID)
	rg.POST("", middlewares.RequireRole("admin", "employee"), ctrl.Create)
	rg.PUT("/:id", middlewares.RequireRole("admin", "employee"), ctrl.Update)
	rg.DELETE("/:id", middlewares.RequireRole("admin"), ctrl.Delete)
}

func (ctrl *FinanceController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("finances")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch finances")
		return
	}
	defer cursor.Close(ctx)

	var finances []models.Finance
	if err = cursor.All(ctx, &finances); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode finances")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, finances)
}

func (ctrl *FinanceController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid finance ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var finance models.Finance
	collection := ctrl.db.Collection("finances")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&finance)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Finance record not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, finance)
}

func (ctrl *FinanceController) Create(c *gin.Context) {
	var finance models.Finance
	if err := c.ShouldBindJSON(&finance); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	finance.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("finances")
	_, err := collection.InsertOne(ctx, finance)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create finance record")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, finance)
}

func (ctrl *FinanceController) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid finance ID")
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

	collection := ctrl.db.Collection("finances")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update finance record")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Finance record not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Finance record updated successfully"})
}

func (ctrl *FinanceController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid finance ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("finances")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete finance record")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Finance record not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Finance record deleted successfully"})
}
