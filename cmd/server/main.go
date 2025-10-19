package main

import (
	"fmt"
	"log"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/handlers"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/routes"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/services"
	"github.com/RaihanurRahman2022/PersonalVault/internal/config"
	"github.com/RaihanurRahman2022/PersonalVault/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type AppConfig struct {
	Config   *config.Config
	Router   *gin.Engine
	Handlers *handlers.Handlers
}

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Faild to initialize application: %v", err)
	}

	log.Println("Started server on port: " + app.Config.Server.Port)
	if err := app.Router.Run(":" + app.Config.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func InitializeApp() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	repos := initializeRepositories(db)

	srvc := initializeService(repos)

	handlers := initializeHandlers(srvc)

	router := configureRouter(handlers, db)

	return &AppConfig{
		Config:   cfg,
		Router:   router,
		Handlers: handlers,
	}, nil
}

func initializeDatabase() (*gorm.DB, error) {
	if err := database.SetupDatabase(); err != nil {
		return nil, fmt.Errorf("failed to setup database: %v", err)
	}

	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initializeRepositories(db *gorm.DB) *repositories.Repositories {
	return repositories.NewRepositories(db)
}

func initializeService(repo *repositories.Repositories) *services.Services {
	return services.NewServices(repo)
}

func initializeHandlers(srvc *services.Services) *handlers.Handlers {
	return handlers.NewHandlers(srvc)
}

func configureRouter(handlers *handlers.Handlers, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	corsConfig := config.PrepareCORSCOnfig()

	router.Use(cors.New(corsConfig))

	routes.SetupRoutes(router, handlers, db)
	return router

}
