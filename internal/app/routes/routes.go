package routes

import (
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/handlers"
	"github.com/RaihanurRahman2022/PersonalVault/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, handlers *handlers.Handlers, db *gorm.DB) {

	setupPublicRoutes(router, handlers.Auth)

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	{
		setupUserRoutes(api, handlers.UserHandler)
		setupDriverRoutes(api, handlers.Driver)
	}
}

// setupPublicRoutes configures public routes that don't require authentication
func setupPublicRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}
}

func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		users.GET("/me", userHandler.GetUserDetails)
	}
}

func setupDriverRoutes(api *gin.RouterGroup, driverHandler *handlers.DriveHandler) {
	driver := api.Group("/drivers")
	{
		driver.GET("/root", driverHandler.GetRootDrivers)
		driver.GET("/list", driverHandler.ListPath)
	}
}
