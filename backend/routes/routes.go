package routes

import (
	"animalsys/config"
	"animalsys/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(rg *gin.RouterGroup, db *mongo.Database, cfg config.Config) {
	controllers.RegisterAuthRoutes(rg.Group("/auth"), db, cfg)
	controllers.RegisterUserRoutes(rg.Group("/users"), db, cfg)
	controllers.RegisterAnimalRoutes(rg.Group("/animals"), db, cfg)
	controllers.RegisterAdoptionRoutes(rg.Group("/adoptions"), db, cfg)
	controllers.RegisterScheduleRoutes(rg.Group("/schedules"), db, cfg)
	controllers.RegisterDocumentRoutes(rg.Group("/documents"), db, cfg)
	controllers.RegisterFinanceRoutes(rg.Group("/finances"), db, cfg)
}
