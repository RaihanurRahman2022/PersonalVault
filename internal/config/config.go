package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Environment string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret       string
	ExpiresInHrs int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	jwtExpiredIn, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expires in hourse time: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "sheikh_enterprise"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "Test_key"),
			ExpiresInHrs: jwtExpiredIn,
		},
		Environment: getEnv("ENV", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func PrepareCORSCOnfig() cors.Config {
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	origins := strings.Split(allowedOrigins, ",")

	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ","),
		AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ","),
		ExposeHeaders:    strings.Split(os.Getenv("CORS_EXPOSE_HEADERS"), ","),
		AllowCredentials: true,
		MaxAge:           300 * time.Second,
	}
}
