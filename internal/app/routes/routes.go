package routes

import (
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/handlers"
	"github.com/RaihanurRahman2022/PersonalVault/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, handlers *handlers.Handlers, db *gorm.DB) {

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	{
		setupUserRoutes(api, handlers.UserHandler)
	}
}

func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		users.GET("/me", userHandler.GetUserDetails)
	}
}
