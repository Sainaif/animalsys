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

type ScheduleController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterScheduleRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &ScheduleController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", ctrl.GetAll)
	rg.GET("/:id", ctrl.GetByID)
	rg.POST("", middlewares.RequireRole("admin", "employee"), ctrl.Create)
	rg.PUT("/:id", ctrl.Update)
	rg.DELETE("/:id", middlewares.RequireRole("admin", "employee"), ctrl.Delete)
}

func (ctrl *ScheduleController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("schedules")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch schedules")
		return
	}
	defer cursor.Close(ctx)

	var schedules []models.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode schedules")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, schedules)
}

func (ctrl *ScheduleController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var schedule models.Schedule
	collection := ctrl.db.Collection("schedules")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&schedule)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Schedule not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, schedule)
}

func (ctrl *ScheduleController) Create(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	schedule.ID = primitive.NewObjectID()
	if schedule.Status == "" {
		schedule.Status = "assigned"
	}
	if schedule.Tasks == nil {
		schedule.Tasks = []string{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("schedules")
	_, err := collection.InsertOne(ctx, schedule)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create schedule")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, schedule)
}

func (ctrl *ScheduleController) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid schedule ID")
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

	collection := ctrl.db.Collection("schedules")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update schedule")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Schedule not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

func (ctrl *ScheduleController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("schedules")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete schedule")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Schedule not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}
