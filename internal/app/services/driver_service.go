package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
)

type DriverService interface {
	GetRoot() ([]entities.RootItems, error)
	ListPath(path string) ([]entities.FileInfo, error)
	Downloadfile(path string) (string, string, error)
}

type DriverServiceImpl struct {
	DriverRepo repositories.DriverRepository
}

func NewDriverService(DriverRepo repositories.DriverRepository) DriverService {
	return &DriverServiceImpl{
		DriverRepo: DriverRepo,
	}
}

func (r *DriverServiceImpl) GetRoot() ([]entities.RootItems, error) {
	log.Println("Service: Getting root directories")
	roots, err := r.DriverRepo.GetRoots()
	if err != nil {
		log.Printf("Service: Error getting root directories: %v", err)
		return nil, fmt.Errorf("failed to get root directories: %w", err)
	}
	if len(roots) == 0 {
		log.Println("Service: No root directories found")
		return []entities.RootItems{}, nil // Return empty slice instead of error
	}

	var rootItems []entities.RootItems
	for _, path := range roots {
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			continue
		}
		name := filepath.Base(path)
		if path == "/" || strings.HasSuffix(path, ":/") {
			name = path
		}

		rootItems = append(rootItems, entities.RootItems{
			Name:     name,
			Path:     path,
			Type:     "directory",
			Size:     0,
			Modified: info.ModTime(),
		})
	}
	log.Printf("Service: Successfully processed %d root items", len(rootItems))
	return rootItems, nil
}

func (r *DriverServiceImpl) ListPath(path string) ([]entities.FileInfo, error) {
	log.Printf("Service: Listing contents of path: %s", path)
	files, err := r.DriverRepo.ListPath(path)
	if err != nil {
		log.Printf("Service: Error listing contents of path: %v", err)
		return nil, fmt.Errorf("failed to list contents of path: %w", err)
	}

	log.Printf("Service: Successfully processed %d files in directory %s", len(files), path)
	return files, nil
}

func (r *DriverServiceImpl) Downloadfile(path string) (string, string, error) {
	log.Printf("Service: Downloading file from path: %s", path)
	absPath, err := r.DriverRepo.Downloadfile(path)
	if err != nil {
		log.Printf("Service: Error downloading file: %v", err)
		return "", "", fmt.Errorf("failed to download file: %w", err)
	}

	filename := filepath.Base(absPath)
	return absPath, filename, nil
}
