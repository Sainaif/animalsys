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

type AdoptionRequest struct {
	AnimalID        string                 `json:"animal_id" binding:"required"`
	ApplicationData map[string]interface{} `json:"application_data" binding:"required"`
}

type AdoptionStatusRequest struct {
	Status string `json:"status" binding:"required"` // approved | rejected
}

func RegisterAdoptionRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	adoptions := db.Collection("adoptions")
	animals := db.Collection("animals")

	rg.GET("/", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cur, err := adoptions.Find(ctx, bson.M{})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.Adoption
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
		var adoption models.Adoption
		err = adoptions.FindOne(ctx, bson.M{"_id": objID}).Decode(&adoption)
		if err != nil {
			utils.Error(c, http.StatusNotFound, "Adoption not found")
			return
		}
		utils.Success(c, adoption)
	})

	rg.POST("/", func(c *gin.Context) {
		var req AdoptionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid animal ID")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var animal models.Animal
		err = animals.FindOne(ctx, bson.M{"_id": animalID, "status": "available"}).Decode(&animal)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Animal not available")
			return
		}
		userID := primitive.NilObjectID
		if uid, ok := c.Get("user_id"); ok {
			userID, _ = primitive.ObjectIDFromHex(uid.(string))
		}
		adoption := models.Adoption{
			AnimalID:        animalID,
			UserID:          userID,
			Status:          "pending",
			ApplicationData: req.ApplicationData,
		}
		res, err := adoptions.InsertOne(ctx, adoption)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		adoption.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, adoption)
	})

	rg.PUT("/:id/status", middlewares.AuthMiddleware(cfg), middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var req AdoptionStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = adoptions.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": req.Status}})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "status updated")
	})

	rg.POST("/:id/contract", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		// Placeholder: just store filename
		file, err := c.FormFile("file")
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "No file uploaded")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// In real app, upload to GridFS or storage and store ObjectID
		_, err = adoptions.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"contract_document_id": file.Filename}})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, "contract uploaded")
	})
}
