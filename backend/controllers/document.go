package controllers

import (
	"animalsys/middlewares"
	"animalsys/models"
	"animalsys/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterDocumentRoutes(rg *gin.RouterGroup, db *mongo.Database) {
	documents := db.Collection("documents")
	os.MkdirAll("./uploads", 0755)

	rg.POST("/", middlewares.AuthMiddleware, func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "No file uploaded")
			return
		}
		f, err := file.Open()
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "File open error")
			return
		}
		defer f.Close()
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		path := filepath.Join("./uploads", filename)
		out, err := os.Create(path)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "File save error")
			return
		}
		defer out.Close()
		size, err := io.Copy(out, f)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "File write error")
			return
		}
		userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
		doc := models.Document{
			Filename:    filename,
			UploadedBy:  userID,
			UploadDate:  time.Now().Format(time.RFC3339),
			Type:        c.PostForm("type"),
			ContentType: file.Header.Get("Content-Type"),
			Size:        size,
		}
		res, err := documents.InsertOne(context.Background(), doc)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		doc.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, doc)
	})

	rg.GET("/", func(c *gin.Context) {
		cur, err := documents.Find(context.Background(), bson.M{})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var docs []models.Document
		if err := cur.All(context.Background(), &docs); err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, docs)
	})

	rg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}
		var doc models.Document
		err = documents.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&doc)
		if err != nil {
			utils.Error(c, http.StatusNotFound, "Document not found")
			return
		}
		c.FileAttachment(filepath.Join("./uploads", doc.Filename), doc.Filename)
	})
}
