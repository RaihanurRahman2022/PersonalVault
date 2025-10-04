package services

import (
	"context"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
)

const previewRangeThreshold = 10 * 1024 * 1024

type DriverService interface {
	GetRoot() ([]entities.RootItems, error)
	ListPath(ctx context.Context, path string) ([]entities.FileInfo, error)
	Downloadfile(path string) (string, string, error)
	CreateFolder(path string) error

	PreviewFile(path string) (*entities.PreviewInfo, error)
	StreamFile(path string) (*entities.PreviewInfo, error)
	UploadFiles(destPath string, files []*multipart.FileHeader, overwrite bool) ([]entities.UploadResult, error)
	UploadFolder(destPath string, files []*multipart.FileHeader, overwrite bool) ([]entities.UploadResult, error)
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
			Name:     name + path,
			Path:     path,
			Type:     "directory",
			Size:     0,
			Modified: info.ModTime(),
		})
	}
	log.Printf("Service: Successfully processed %d root items", len(rootItems))
	return rootItems, nil
}

// Example of using context to cancel the operation
func (r *DriverServiceImpl) ListPath(ctx context.Context, path string) ([]entities.FileInfo, error) {
	log.Printf("Service: Listing contents of path: %s", path)

	// Check if context is cancelled before proceeding
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	files, err := r.DriverRepo.ListPath(ctx, path)
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

func (r *DriverServiceImpl) CreateFolder(path string) error {
	log.Printf("Service: Creating folder at path: %s", path)
	err := r.DriverRepo.CreateFolder(path)
	if err != nil {
		log.Printf("Service: Error creating folder: %v", err)
		return fmt.Errorf("failed to create folder: %w", err)
	}
	log.Printf("Service: Successfully created folder at path: %s", path)
	return nil
}

func (r *DriverServiceImpl) PreviewFile(path string) (*entities.PreviewInfo, error) {
	log.Printf("Service: Previewing file at path: %s", path)

	file, fileInfo, absPath, err := r.DriverRepo.OpenFile(path)

	if err != nil {
		log.Printf("Service: Error opening file: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	mimeType := detectMimeType(absPath)

	shouldUseRange := fileInfo.Size() >= previewRangeThreshold

	return &entities.PreviewInfo{
		File:           file,
		Info:           fileInfo,
		AbsPath:        absPath,
		MimeType:       mimeType,
		ShouldUseRange: shouldUseRange,
	}, nil
}

func (r *DriverServiceImpl) StreamFile(path string) (*entities.PreviewInfo, error) {
	log.Printf("Service: Streaming file at path: %s", path)

	file, fileInfo, absPath, err := r.DriverRepo.OpenFile(path)

	if err != nil {
		log.Printf("Service: Error opening file: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	mimeType := detectMimeType(absPath)
	return &entities.PreviewInfo{
		File:           file,
		Info:           fileInfo,
		AbsPath:        absPath,
		MimeType:       mimeType,
		ShouldUseRange: true,
	}, nil
}

func detectMimeType(absPath string) string {
	ext := strings.ToLower(filepath.Ext(absPath))
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	return "applicaton/octet-stream"
}

// UploadFiles uploads a list of individual files to the destination path using concurrent goroutines.
//
// ⚠️ Improvements / Considerations:
// 1. Currently it launches one goroutine per file without limit. For a large number of files,
//    this may consume a lot of memory and CPU. Consider using a worker pool or limiting concurrency.
// 2. No context is used. If the client cancels the request (e.g., closes the browser tab), the uploads
//    will continue in the background. Using `context.Context` would allow graceful cancellation.
// 3. Logging is minimal. Consider adding per-file logging for better observability in production.
// 4. Error handling is per-file and collected in the slice. This is fine, but you could also use
//    `errgroup` if you want to stop all uploads on first error.
//
// Important parts for new developers:
// - `wg.Wait()` ensures the main goroutine waits for all upload goroutines to finish.
// - Writing results into `result[index]` is safe because each index is unique per goroutine.

func (r *DriverServiceImpl) UploadFiles(destPath string, files []*multipart.FileHeader, overwrite bool) ([]entities.UploadResult, error) {
	log.Printf("Service: Uploading %d files to path: %s", len(files), destPath)

	if err := r.DriverRepo.EnsureDirExists(destPath); err != nil {
		return nil, fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	var wg sync.WaitGroup
	result := make([]entities.UploadResult, len(files))

	for i, file := range files {
		wg.Add(1)
		go func(index int, fh *multipart.FileHeader) {
			defer wg.Done()

			res := entities.UploadResult{Name: fh.Filename}
			dst := filepath.Join(destPath, fh.Filename)

			written, err := r.DriverRepo.SaveUploadedFile(fh, dst, overwrite)
			if err != nil {
				res.Error = err.Error()
			} else {
				res.Path = dst
				res.Size = written
			}
			// Writing directly to the slice by index is safe because each goroutine has a unique index
			result[index] = res

		}(i, file)
	}

	wg.Wait() // wait until all files processed

	allFailed := true
	for _, res := range result {
		if res.Error == "" {
			allFailed = false
			break
		}
	}

	if allFailed {
		return nil, fmt.Errorf("all files failed to upload")
	}

	log.Printf("Service: Successfully uploaded %d files to path: %s", len(result), destPath)

	return result, nil
}

// UploadFolder uploads a folder containing multiple files using a worker pool pattern.
//
// ⚠️ Improvements / Considerations:
// 1. Currently, the jobs channel is buffered with `len(files)`. For huge folders (e.g., 100k files),
//    this may consume a lot of memory. Consider a smaller buffer and push jobs gradually.
// 2. No context is used. If the client disconnects, the uploads continue. Adding `context.Context`
//    allows graceful shutdown of ongoing uploads.
// 3. Worker count is hardcoded to 5. In production, this should be configurable based on CPU/network IO.
// 4. Results are collected in the order they complete, not necessarily the order of `files`. If order matters,
//    additional logic is needed.
//
// Important parts for new developers:
// - Worker pool: fixed number of workers consume jobs from the channel.
// - Each worker ensures parent directories exist before saving files.
// - Channels `jobs` and `result` are used for synchronization between main goroutine and workers.

func (r *DriverServiceImpl) UploadFolder(destPath string, files []*multipart.FileHeader, overwrite bool) ([]entities.UploadResult, error) {
	if err := r.DriverRepo.EnsureDirExists(destPath); err != nil {
		return nil, fmt.Errorf("invalid destination: %w", err)
	}

	const workerCount = 5

	jobs := make(chan *multipart.FileHeader, len(files))   // ⚠️ Consider smaller buffer for huge folders
	result := make(chan entities.UploadResult, len(files)) // ⚠️ Consider smaller buffer for huge folders

	for range workerCount {

		go func() {
			for fh := range jobs {
				res := entities.UploadResult{Name: fh.Filename}

				// Debug: log the filename received from frontend
				log.Printf("Backend received file: %s", fh.Filename)

				// For folder uploads, the filename might contain relative path
				// e.g., "folder/subfolder/file.txt"
				fullPath := filepath.Join(destPath, fh.Filename)

				log.Printf("Backend received file: %s", fullPath)

				// Ensure parent directories exist
				parentDir := filepath.Dir(fullPath)
				if err := r.DriverRepo.EnsureDirExists(parentDir); err != nil {
					res.Error = fmt.Sprintf("failed to create directory: %v", err)
					result <- res
					continue
				}

				written, err := r.DriverRepo.SaveUploadedFile(fh, fullPath, overwrite)
				if err != nil {
					res.Error = err.Error()
				} else {
					res.Path = fullPath
					res.Size = written
				}
				result <- res
			}
		}()

	}

	// Enqueue all jobs
	for _, fh := range files {
		jobs <- fh
	}

	close(jobs) // Signal workers no more jobs will be sent

	var uploadedResult []entities.UploadResult
	for range files {
		uploadedResult = append(uploadedResult, <-result)
	}
	// Check if all failed
	allFailed := true
	for _, r := range uploadedResult {
		if r.Error == "" {
			allFailed = false
			break
		}
	}
	if allFailed {
		return uploadedResult, fmt.Errorf("all uploads failed")
	}
	return uploadedResult, nil
}
