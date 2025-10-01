package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func (h *DriveHandler) CreateFolder(c *gin.Context) {
	var req entities.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body " + err.Error()})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	log.Printf("Handler: Creating folder at path: %s", req.Path)
	err := h.DriverService.CreateFolder(req.Path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder created successfully"})
}

func (h *DriveHandler) PreviewFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	log.Printf("Handler: Previewing file at path: %s", path)
	previewInfo, err := h.DriverService.PreviewFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer previewInfo.File.Close()

	c.Header("Content-Type", previewInfo.MimeType)
	c.Header("Cache-Control", "public, max-age=3600")

	if !previewInfo.ShouldUseRange {
		h.returnWholeFile(c, previewInfo.File, previewInfo.Info.Size())
		return
	}

	h.serveWithRange(c, previewInfo.File, previewInfo.Info)
}

func (h *DriveHandler) StreamFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	streamInfo, err := h.DriverService.StreamFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer streamInfo.File.Close()

	c.Header("Content-Type", streamInfo.MimeType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "public, max-age=3600")

	h.serveWithRange(c, streamInfo.File, streamInfo.Info)
}

func (h *DriveHandler) serveWithRange(c *gin.Context, file *os.File, info os.FileInfo) {
	fileSize := info.Size()
	rangeHeader := c.GetHeader("Range")

	if rangeHeader == "" {
		h.returnWholeFile(c, file, fileSize)
		return
	}

	start, end, err := h.parseRangeHeader(c, rangeHeader, file, fileSize)
	if err != nil {
		return
	}

	contentLength := end - start + 1
	c.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileSize, 10))
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Status(http.StatusPartialContent)

	_, err = io.CopyN(c.Writer, file, contentLength)
	if err != nil {
		log.Printf("Error serving file: %v", err)
		return
	}
}

func (h *DriveHandler) returnWholeFile(c *gin.Context, file *os.File, fileSize int64) {
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	_, err := io.Copy(c.Writer, file)
	if err != nil {
		log.Printf("Error serving file: %v", err)
		return
	}
}

func (h *DriveHandler) parseRangeHeader(c *gin.Context, rangeHeader string, file *os.File, fileSize int64) (int64, int64, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return 0, 0, fmt.Errorf("invalid range header")
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return 0, 0, fmt.Errorf("invalid range header")
	}

	var start, end int64
	var err error

	if parts[0] == "" {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil || start < 0 {
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return 0, 0, fmt.Errorf("invalid range header")
		}
	}

	if parts[1] == "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end >= fileSize {
			end = fileSize - 1
		}
	} else {
		end = fileSize - 1
	}

	if start > end || start >= fileSize {
		c.Header("Content-Range", "bytes */"+strconv.FormatInt(fileSize, 10))
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return 0, 0, fmt.Errorf("invalid range header")
	}

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return 0, 0, fmt.Errorf("failed to seek file: %v", err)
	}

	return start, end, nil
}
