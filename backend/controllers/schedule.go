package controllers

import (
	"animalsys/middlewares"
	"animalsys/models"
	"animalsys/utils"
	"context"
	"net/http"
	"time"

	"animalsys/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScheduleRequest struct {
	EmployeeID string   `json:"employee_id" binding:"required"`
	ShiftDate  string   `json:"shift_date" binding:"required"`
	ShiftTime  string   `json:"shift_time" binding:"required"`
	Tasks      []string `json:"tasks" binding:"required"`
}

type ScheduleStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type SwapRequest struct {
	TargetEmployeeID string `json:"target_employee_id" binding:"required"`
}

type AbsenceRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func RegisterScheduleRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	schedules := db.Collection("schedules")

	rg.GET("/", middlewares.AuthMiddleware(cfg), func(c *gin.Context) {
		role := c.GetString("role")
		userID := c.GetString("user_id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{}
		if role == "volunteer" {
			uid, _ := primitive.ObjectIDFromHex(userID)
			filter["employee_id"] = uid
		}
		cur, err := schedules.Find(ctx, filter)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.Schedule
		if err := cur.All(ctx, &result); err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, result)
	})

	rg.POST("/", middlewares.AuthMiddleware(cfg), middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		var req ScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		employeeID, err := primitive.ObjectIDFromHex(req.EmployeeID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid employee ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		schedule := models.Schedule{
			EmployeeID: employeeID,
			ShiftDate:  req.ShiftDate,
			ShiftTime:  req.ShiftTime,
			Tasks:      req.Tasks,
			Status:     "assigned",
		}
		res, err := schedules.InsertOne(ctx, schedule)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		schedule.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, schedule)
	})

	rg.PUT("/:id/swap", middlewares.AuthMiddleware(cfg), middlewares.RBACMiddleware("volunteer"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var req SwapRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		targetID, err := primitive.ObjectIDFromHex(req.TargetEmployeeID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid target ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = schedules.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": "swap_requested", "target_employee_id": targetID}})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "swap requested")
	})

	rg.PUT("/:id/absence", middlewares.AuthMiddleware(cfg), middlewares.RBACMiddleware("volunteer"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var req AbsenceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = schedules.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": "absent_requested", "absence_reason": req.Reason}})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "absence requested")
	})

	rg.PUT("/:id/status", middlewares.AuthMiddleware(cfg), middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var req ScheduleStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = schedules.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": req.Status}})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "status updated")
	})
}
