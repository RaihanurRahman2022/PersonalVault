package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
)

type DriverService interface {
	GetRoot() ([]entities.RootItems, error)
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
	roots, err := r.DriverRepo.GetRoots()
	if err != nil {
		return nil, err
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("no valid root directories found")
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
	return rootItems, nil
}
