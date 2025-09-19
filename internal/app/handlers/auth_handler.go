package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	//Read the Raw body
	body, _ := c.GetRawData()
	fmt.Println("Raw login request body:", string(body))

	var jsonData map[string]any
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
	}

	fmt.Printf("Parsed JSON: %+v\n", jsonData)

	req := entities.LoginRequest{
		Username: getString(jsonData, "username"),
		Password: getString(jsonData, "password"),
	}

	accessToken, refreshToken, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		if err == services.ErrUserInactive {
			c.JSON(http.StatusForbidden, gin.H{"error": "user account is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {

	body, _ := c.GetRawData()
	fmt.Println("Raw request body: ", string(body))

	var jsonData map[string]any
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Printf("JSON parsing error %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON format",
		})
		return
	}

	fmt.Printf("Parsed JSON %+v\n", jsonData)

	req := entities.RegisterRequest{
		Username:  getString(jsonData, "username"),
		FirstName: getString(jsonData, "first_name"),
		LastName:  getString(jsonData, "last_name"),
		Password:  getString(jsonData, "password"),
	}

	if err := h.authService.Register(&req); err != nil {
		fmt.Printf("Registration error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register the user",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

// Helper to get string value from map
func getString(data map[string]any, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
