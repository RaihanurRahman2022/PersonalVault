package handlers

import (
	"log"
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

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"Data":    roots,
		"message": "fetch all drivers successfully",
	})
}

func (h *DriveHandler) ListPath(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	log.Printf("Handler: Listing contents of path: %s", path)
	files, err := h.DriverService.ListPath(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list path " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"Data":    files,
		"message": "fetch all files successfully",
	})
}

func (h *DriveHandler) Downloadfile(c *gin.Context) {
	var req entities.DownloadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body " + err.Error()})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	log.Printf("Handler: Downloading file from path: %s", req.Path)
	absPath, filename, err := h.DriverService.Downloadfile(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file " + err.Error()})
		return
	}

	c.FileAttachment(absPath, filename)
}
