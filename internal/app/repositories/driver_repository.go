package repositories

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"golang.org/x/sys/windows"
)

type DriverRepository interface {
	GetRoots() ([]string, error)
	ListPath(path string) ([]entities.FileInfo, error)
	Downloadfile(path string) (string, error)
}

type DriverRepositoryImpl struct {
}

func NewDriverRepository() DriverRepository {
	return &DriverRepositoryImpl{}
}

func (r *DriverRepositoryImpl) GetRoots() ([]string, error) {
	var roots []string
	log.Printf("Getting root directories for OS: %s", runtime.GOOS)

	if runtime.GOOS == "windows" {
		Drivers, err := getWindowsDrivers()
		if err != nil {
			log.Printf("Error getting Windows drivers: %v", err)
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

	log.Printf("Found %d unique root directories", len(result))
	return result, nil
}

func getWindowsDrivers() ([]string, error) {
	var drivers []string

	// Each drive string is like "C:\\" and strings are concatenated with \0 separators, ending with \0\0.
	// The API expects the buffer size in UTF-16 code units (not bytes).
	buf := make([]uint16, 256)
	n, err := windows.GetLogicalDriveStrings(uint32(len(buf)), &buf[0])
	if err != nil {
		return nil, err
	}

	if n == 0 {
		// Fallback to C:/ if API returns nothing
		return []string{"C:/"}, nil
	}

	// Ensure we only parse up to n code units returned by the API
	u := buf[:n]

	// Parse sequences separated by 0 (NUL). There is a trailing double NUL; ignore empties.
	start := 0
	for i, v := range u {
		if v == 0 {
			if i > start {
				s := windows.UTF16ToString(u[start:i])
				if s != "" {
					s = strings.ReplaceAll(s, "\\", "/")
					drivers = append(drivers, s)
				}
			}
			start = i + 1
		}
	}

	log.Printf("Drivers found (parsed): %v", drivers)
	if len(drivers) == 0 {
		return []string{"C:/"}, nil
	}
	return drivers, nil
}

func (r *DriverRepositoryImpl) ListPath(path string) ([]entities.FileInfo, error) {
	log.Printf("Repository: Listing contents of path: %s", path)

	if !isSafePath(path) {
		return nil, fmt.Errorf("access to path %s is not allowed", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path does not exist or is not accessible: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fileinfos []entities.FileInfo

	for _, entry := range entries {
		if shouldSkipFile(entry) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("repository: failed to get file info: %w", err)
		}
		fileType := "file"
		if info.IsDir() {
			fileType = "folder"
		}

		fullPath := filepath.Join(path, entry.Name())
		fileinfo := entities.FileInfo{
			Name:     entry.Name(),
			Path:     fullPath,
			Type:     fileType,
			Size:     info.Size(),
			Modified: info.ModTime(),
		}
		fileinfos = append(fileinfos, fileinfo)
	}

	log.Printf("Repository: Successfully processed %d files in directory %s", len(fileinfos), path)
	return fileinfos, nil
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

func shouldSkipFile(entry os.DirEntry) bool {
	filename := entry.Name()

	if strings.HasPrefix(filename, ".") {
		return true
	}

	if strings.HasPrefix(filename, "$") {
		return true
	}

	systemFolders := []string{
		"System Volume Information",
		"Recovery",
		"Windows",
		"Program Files",
		"Program Files (x86)",
		"ProgramData",
		"Boot",
		"EFI",
	}

	for _, folder := range systemFolders {
		if strings.HasPrefix(filename, folder) {
			return true
		}
	}

	return false
}

func (r *DriverRepositoryImpl) Downloadfile(path string) (string, error) {
	if !isSafePath(path) {
		return "", fmt.Errorf("access to path %s is not allowed", path)
	}
	info, err := os.Stat(path)

	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", fmt.Errorf("folder or directory is not downloadable")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
