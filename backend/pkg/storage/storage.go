package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StorageService handles file uploads and management
type StorageService struct {
	basePath   string
	baseURL    string
	maxFileSize int64 // in bytes
}

// NewStorageService creates a new storage service
func NewStorageService(basePath, baseURL string, maxFileSize int64) *StorageService {
	return &StorageService{
		basePath:    basePath,
		baseURL:     baseURL,
		maxFileSize: maxFileSize,
	}
}

// UploadImage uploads an image file and returns the URL
func (s *StorageService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	// Validate file size
	if file.Size > s.maxFileSize {
		return "", errors.NewBadRequest(fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", s.maxFileSize))
	}

	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return "", errors.NewBadRequest("invalid file type. Allowed types: jpg, jpeg, png, gif, webp")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", errors.Wrap(err, 500, "failed to open uploaded file")
	}
	defer src.Close()

	// Generate unique filename
	filename := s.generateFilename(file.Filename)

	// Create folder path
	folderPath := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", errors.Wrap(err, 500, "failed to create upload directory")
	}

	// Create destination file
	destPath := filepath.Join(folderPath, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", errors.Wrap(err, 500, "failed to create destination file")
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.Wrap(err, 500, "failed to save file")
	}

	// Generate URL
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, folder, filename)

	return url, nil
}

// UploadMultipleImages uploads multiple images
func (s *StorageService) UploadMultipleImages(ctx context.Context, files []*multipart.FileHeader, folder string) ([]string, error) {
	urls := make([]string, 0, len(files))

	for _, file := range files {
		url, err := s.UploadImage(ctx, file, folder)
		if err != nil {
			// Clean up previously uploaded files on error
			for _, uploadedURL := range urls {
				_ = s.DeleteFile(ctx, uploadedURL)
			}
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

// DeleteFile deletes a file by URL
func (s *StorageService) DeleteFile(ctx context.Context, url string) error {
	// Extract file path from URL
	relativePath := strings.TrimPrefix(url, s.baseURL+"/")
	filePath := filepath.Join(s.basePath, relativePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.ErrNotFound
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return errors.Wrap(err, 500, "failed to delete file")
	}

	return nil
}

// DeleteMultipleFiles deletes multiple files
func (s *StorageService) DeleteMultipleFiles(ctx context.Context, urls []string) error {
	for _, url := range urls {
		_ = s.DeleteFile(ctx, url) // Continue on error
	}
	return nil
}

// isValidImageType checks if the file extension is valid for images
func (s *StorageService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}

	return false
}

// generateFilename generates a unique filename
func (s *StorageService) generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UnixNano()
	randomID := primitive.NewObjectID().Hex()

	return fmt.Sprintf("%d_%s%s", timestamp, randomID, ext)
}
