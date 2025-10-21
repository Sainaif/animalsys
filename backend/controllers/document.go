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

type DocumentController struct {
	db  *mongo.Database
	cfg config.Config
}

func RegisterDocumentRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	ctrl := &DocumentController{db: db, cfg: cfg}
	rg.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	rg.GET("", ctrl.GetAll)
	rg.GET("/:id", ctrl.GetByID)
	rg.POST("", middlewares.RequireRole("admin", "employee"), ctrl.Create)
	rg.DELETE("/:id", middlewares.RequireRole("admin", "employee"), ctrl.Delete)
}

func (ctrl *DocumentController) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("documents")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch documents")
		return
	}
	defer cursor.Close(ctx)

	var documents []models.Document
	if err = cursor.All(ctx, &documents); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode documents")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, documents)
}

func (ctrl *DocumentController) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid document ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var document models.Document
	collection := ctrl.db.Collection("documents")
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Document not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, document)
}

func (ctrl *DocumentController) Create(c *gin.Context) {
	var document models.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	document.ID = primitive.NewObjectID()
	document.UploadDate = time.Now().Format(time.RFC3339)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("documents")
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create document")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, document)
}

func (ctrl *DocumentController) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid document ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := ctrl.db.Collection("documents")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete document")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Document not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
