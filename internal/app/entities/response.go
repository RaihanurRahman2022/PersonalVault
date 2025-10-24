package entities

import (
	"os"
	"time"
)

// RootItems represents a root directory item
type RootItems struct {
	Name     string    `json:"name" example:"C:"`
	Path     string    `json:"path" example:"C:\\"`
	Type     string    `json:"type" example:"drive"`
	Size     int64     `json:"size" example:"1073741824000"`
	Modified time.Time `json:"modified" example:"2024-01-15T10:30:00Z"`
}

// FileInfo represents file/folder information
type FileInfo struct {
	Name     string    `json:"name" example:"document.pdf"`
	Path     string    `json:"path" example:"/documents/document.pdf"`
	Type     string    `json:"type" example:"file"`
	Size     int64     `json:"size" example:"1024000"`
	Modified time.Time `json:"modified" example:"2024-01-15T10:30:00Z"`
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	Name  string `json:"name" example:"document.pdf"`
	Path  string `json:"path,omitempty" example:"/documents/document.pdf"`
	Size  int64  `json:"size,omitempty" example:"1024000"`
	Error string `json:"error,omitempty" example:"File already exists"`
}

// PreviewInfo represents the information needed for file preview
type PreviewInfo struct {
	File           *os.File    `json:"file" example:"*os.File"`
	Info           os.FileInfo `json:"info" example:"os.FileInfo"`
	AbsPath        string      `json:"abs_path" example:"/documents/document.pdf"`
	MimeType       string      `json:"mime_type" example:"application/pdf"`
	ShouldUseRange bool        `json:"should_use_range" example:"true"`
}
