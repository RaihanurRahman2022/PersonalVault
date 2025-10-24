package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

// GetRootDrivers godoc
// @Summary      Get root drivers
// @Description  Get all root drivers
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]any "Root drivers"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/root [get]
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

// ListPath godoc
// @Summary      List path contents
// @Description  List all files and directories in a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        path query string true "Path to list contents"
// @Success      200 {object} map[string]any "Path contents"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/list [get]
func (h *DriveHandler) ListPath(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	log.Printf("Handler: Listing contents of path: %s", path)
	files, err := h.DriverService.ListPath(ctx, path)
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

// Downloadfile godoc
// @Summary      Download file
// @Description  Download a file from a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        request body entities.DownloadRequest true "Download request"
// @Success      200 {object} map[string]string "File downloaded successfully"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/download [post]
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

// CreateFolder godoc
// @Summary      Create folder
// @Description  Create a new folder at a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        request body entities.CreateFolderRequest true "Create folder request"
// @Success      200 {object} map[string]string "Folder created successfully"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/create-folder [post]
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

// PreviewFile godoc
// @Summary      Preview file
// @Description  Preview a file at a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        path query string true "Path to preview file"
// @Success      200 {object} map[string]string "File previewed successfully"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/preview [get]
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

// StreamFile godoc
// @Summary      Stream file
// @Description  Stream a file at a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        path query string true "Path to stream file"
// @Success      200 {object} map[string]string "File streamed successfully"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/stream [get]
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

	buffer := make([]byte, 32*1024)
	_, err = io.CopyBuffer(c.Writer, file, buffer)
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

func extractFilenameFromHeader(header map[string][]string) string {
	if contentDisposition, exists := header["Content-Disposition"]; exists && len(contentDisposition) > 0 {
		// Parse Content-Disposition header like: form-data; name="files"; filename="Handle/handle.exe"
		disposition := contentDisposition[0]
		log.Printf("Handler: Parsing disposition: %s", disposition)
		if strings.Contains(disposition, "filename=") {
			// Extract filename from the header
			parts := strings.Split(disposition, "filename=")
			log.Printf("Handler: Split parts: %v", parts)
			if len(parts) > 1 {
				filename := strings.Trim(parts[1], `"`)
				log.Printf("Handler: Extracted filename: %s", filename)
				return filename
			}
		}
	}
	return ""
}

// UploadFiles godoc
// @Summary      Upload files
// @Description  Upload files to a given path
// @Tags         Drive
// @Accept       json
// @Produce      json
// @Param        path query string true "Path to upload files"
// @Param        upload_type formData string true "Upload type (files or folder)"
// @Param        overwrite formData string true "Overwrite files (true or false)"
// @Success      200 {object} map[string]any "Upload results"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /drive/upload-files [post]
func (h *DriveHandler) UploadFiles(c *gin.Context) {
	dstPath := c.Query("path")
	if dstPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "destination path is required"})
		return
	}

	uploadType := c.DefaultPostForm("upload_type", "files")
	overwrite := strings.ToLower(c.DefaultPostForm("overwrite", "false")) == "true"

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form: " + err.Error()})
		return
	}

	var results []entities.UploadResult

	if uploadType == "folder" {
		results, err = h.handleFolderUpload(dstPath, form, overwrite)
	} else {
		results, err = h.handleFileUpload(c, dstPath, form, overwrite)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"message": "files uploaded successfully",
	})

}

func (h *DriveHandler) handleFileUpload(c *gin.Context, destPath string, form *multipart.Form, overwrite bool) ([]entities.UploadResult, error) {
	files := form.File["files"]
	if len(files) == 0 {
		if f, err := c.FormFile("file"); err == nil {
			files = []*multipart.FileHeader{f}
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// Extract correct filename from Content-Disposition header for file uploads too
	for i, file := range files {
		if file.Header != nil {
			extractedFilename := extractFilenameFromHeader(file.Header)
			if extractedFilename != "" {
				log.Printf("Handler: File upload %d - Extracted filename: %s", i, extractedFilename)
				file.Filename = extractedFilename
			}
		}
	}

	return h.DriverService.UploadFiles(destPath, files, overwrite)
}

func (h *DriveHandler) handleFolderUpload(destPath string, form *multipart.Form, overwrite bool) ([]entities.UploadResult, error) {
	files := form.File["files"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// Debug: log the raw filename from multipart form
	for i, file := range files {
		log.Printf("Handler: Raw file %d - Filename: %s, Size: %d", i, file.Filename, file.Size)
		// Log the header to see what's in the Content-Disposition
		if file.Header != nil {
			// Extract filename from Content-Disposition header
			extractedFilename := extractFilenameFromHeader(file.Header)
			if extractedFilename != "" {
				log.Printf("Handler: Extracted filename from header: %s", extractedFilename)
				// Update the filename to preserve the folder structure
				file.Filename = extractedFilename
				log.Printf("Handler: Updated filename to: %s", file.Filename)
			} else {
				log.Printf("Handler: No filename extracted for file %d", i)
			}
		}
	}

	return h.DriverService.UploadFolder(destPath, files, overwrite)
}
