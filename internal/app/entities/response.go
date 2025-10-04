package entities

import (
	"os"
	"time"
)

type RootItems struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Type     string    `json:"type"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Type     string    `json:"type"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

type PreviewInfo struct {
	File           *os.File
	Info           os.FileInfo
	AbsPath        string
	MimeType       string
	ShouldUseRange bool
}

type UploadResult struct {
	Name  string `json:"name"`
	Path  string `json:"path,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Error string `json:"error,omitempty"`
}
