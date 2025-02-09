package service

import (
	"context"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// LocalStorageService implements StorageService using local filesystem
type LocalStorageService struct {
	config *config.Config
}

// NewLocalStorageService creates a new LocalStorageService
func NewLocalStorageService(config *config.Config) StorageService {
	return &LocalStorageService{
		config: config,
	}
}

func (s *LocalStorageService) DownloadFile(ctx context.Context, url string) (string, error) {
	logger.Debug("downloading file", "url", url)

	// Create temporary file
	tmpFile, err := os.CreateTemp(s.config.Storage.TempDirectory, "input-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Copy file contents
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return tmpFile.Name(), nil
}

func (s *LocalStorageService) UploadFile(ctx context.Context, localPath string, objectKey string, userID uint) (string, error) {
	// For local storage, we'll just copy the file to a permanent location
	userPath := fmt.Sprintf("user_%d", userID)
	destPath := filepath.Join(s.config.Storage.TempDirectory, "uploads", userPath, objectKey)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	srcFile, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return URL
	return fmt.Sprintf("%s/%s/%s", s.config.Storage.TempDirectory, userPath, objectKey), nil
}

func (s *LocalStorageService) DeleteFile(ctx context.Context, localPath string) error {
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
