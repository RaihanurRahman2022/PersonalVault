package handlers

import (
	"net/http"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/services"
	"github.com/gin-gonic/gin"
)

type DriveHandler struct {
	DriverService services.DriverService
}

func NewDriverHandler(srvc services.DriverService) *DriveHandler {
	return &DriveHandler{
		DriverService: srvc,
	}
}

func (h *DriveHandler) GetRootDrivers(c *gin.Context) {
	roots, err := h.DriverService.GetRoot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roots " + err.Error()})
		return
	}

	if len(roots) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Data":    []entities.RootItems{},
			"message": "No root directories found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Data":    roots,
		"message": "fetch all drivers successfully",
	})
}
