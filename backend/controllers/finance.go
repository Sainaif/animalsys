package controllers

import (
	"animalsys/middlewares"
	"animalsys/models"
	"animalsys/utils"
	"context"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FinanceRequest struct {
	Date        string  `json:"date" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

func RegisterFinanceRoutes(rg *gin.RouterGroup, db *mongo.Database) {
	finances := db.Collection("finances")

	rg.GET("/", middlewares.AuthMiddleware, func(c *gin.Context) {
		cur, err := finances.Find(context.Background(), bson.M{})
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.Finance
		if err := cur.All(context.Background(), &result); err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		utils.Success(c, result)
	})

	rg.POST("/", middlewares.AuthMiddleware, middlewares.RBACMiddleware("admin", "employee"), func(c *gin.Context) {
		var req FinanceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid input")
			return
		}
		userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
		finance := models.Finance{
			Date:        req.Date,
			Type:        req.Type,
			Amount:      req.Amount,
			Description: req.Description,
			Category:    req.Category,
			CreatedBy:   userID,
		}
		res, err := finances.InsertOne(context.Background(), finance)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "DB error")
			return
		}
		finance.ID = res.InsertedID.(primitive.ObjectID)
		utils.Success(c, finance)
	})

	rg.GET("/report/csv", middlewares.AuthMiddleware, func(c *gin.Context) {
		cur, err := finances.Find(context.Background(), bson.M{})
		if err != nil {
			c.String(http.StatusInternalServerError, "DB error")
			return
		}
		var result []models.Finance
		if err := cur.All(context.Background(), &result); err != nil {
			c.String(http.StatusInternalServerError, "DB error")
			return
		}
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=finances.csv")
		w := csv.NewWriter(c.Writer)
		w.Write([]string{"Date", "Type", "Amount", "Description", "Category"})
		for _, f := range result {
			w.Write([]string{f.Date, f.Type, strconv.FormatFloat(f.Amount, 'f', 2, 64), f.Description, f.Category})
		}
		w.Flush()
	})
}
