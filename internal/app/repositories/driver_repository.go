package repositories

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
	"gorm.io/gorm"
)

type DriverRepository interface {
	GetRoots() ([]string, error)
}

type DriverRepositoryImpl struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &DriverRepositoryImpl{
		db: db,
	}
}

func (r *DriverRepositoryImpl) GetRoots() ([]string, error) {
	var roots []string

	if runtime.GOOS == "windows" {
		Drivers, err := getWindowsDrivers()
		if err != nil {
			return nil, err
		}

		for _, d := range Drivers {
			absPath, err := filepath.Abs(filepath.Clean(d))
			if err == nil && isSafePath(absPath) {
				roots = append(roots, absPath)
			}
		}
	} else {
		roots = append(roots, "/")
		if home, err := os.UserHomeDir(); err == nil {
			absHome, err := filepath.Abs(filepath.Clean(home))
			if err == nil && isSafePath(absHome) {
				roots = append(roots, absHome)
			}
		}

		mountDirs := []string{"/mnt", "/media"}
		for _, dir := range mountDirs {
			if entries, err := os.ReadDir(dir); err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						absPath, err := filepath.Abs(filepath.Clean(filepath.Join(dir, entry.Name())))
						if err == nil && isSafePath(absPath) {
							roots = append(roots, absPath)
						}
					}
				}
			}
		}
	}

	uniqueRoots := make(map[string]bool)
	var result []string
	for _, root := range roots {
		if !uniqueRoots[root] {
			uniqueRoots[root] = true
			result = append(result, root)
		}
	}

	return result, nil
}

func getWindowsDrivers() ([]string, error) {
	var Drivers []string
	buffer := make([]uint16, 1024)
	n, err := windows.GetLogicalDriveStrings(uint32(len(buffer)/2), &buffer[0])
	if err != nil {
		return nil, err
	}

	DriverStrings := windows.UTF16ToString(buffer[:n])
	if DriverStrings == "" {
		return []string{"C:/"}, nil
	}

	for _, d := range strings.Split(DriverStrings, "\x00") {
		if d != "" {
			d = strings.ReplaceAll(d, "\\", "/")
			Drivers = append(Drivers, d)
		}
	}

	return Drivers, nil
}

// isSafePath ensures the path is safe (e.g., not accessing sensitive system dirs)
func isSafePath(path string) bool {
	// Block sensitive paths (customize as needed)
	sensitivePaths := []string{
		"/etc",
		"/var",
		"/root",
		"/proc",
		"/sys",
	}
	if runtime.GOOS == "windows" {
		sensitivePaths = []string{
			"C:/Windows",
			"C:/Program Files",
			"C:/Program Files (x86)",
		}
	}
	for _, sp := range sensitivePaths {
		if strings.HasPrefix(strings.ToLower(path), strings.ToLower(sp)) {
			return false
		}
	}
	return true
}
