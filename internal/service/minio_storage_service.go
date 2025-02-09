package service

import (
	"context"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorageService implements StorageService using MinIO
type MinioStorageService struct {
	config *config.Config
	client *minio.Client
}

// NewMinioStorageService creates a new MinioStorageService
func NewMinioStorageService(config *config.Config) (StorageService, error) {
	// Initialize MinIO client
	endpoint := fmt.Sprintf("%s:%s", config.Storage.MinioEndpoint, config.Storage.MinioPort)
	logger.Info("initializing MinIO client", "endpoint", endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Storage.MinioAccessKey, config.Storage.MinioSecretKey, ""),
		Secure: config.Storage.MinioUseSSL,
		Region: config.Storage.MinioRegion,
	})

	if err != nil {
		logger.Error("failed to create MinIO client", "error", err)
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create bucket if it doesn't exist
	exists, err := client.BucketExists(context.Background(), config.Storage.MinioBucketName)
	if err != nil {
		logger.Error("failed to check bucket existence",
			"bucket", config.Storage.MinioBucketName,
			"error", err)
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		logger.Info("creating bucket", "bucket", config.Storage.MinioBucketName)
		err = client.MakeBucket(context.Background(), config.Storage.MinioBucketName, minio.MakeBucketOptions{
			Region: config.Storage.MinioRegion,
		})
		if err != nil {
			logger.Error("failed to create bucket",
				"bucket", config.Storage.MinioBucketName,
				"error", err)
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		logger.Info("bucket created successfully", "bucket", config.Storage.MinioBucketName)
	}

	// Set bucket policy to public
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + config.Storage.MinioBucketName + `/*"]
			}
		]
	}`

	err = client.SetBucketPolicy(context.Background(), config.Storage.MinioBucketName, policy)
	if err != nil {
		logger.Error("failed to set bucket policy", "error", err)
		return nil, fmt.Errorf("failed to set bucket policy: %w", err)
	}

	logger.Info("MinIO storage service initialized successfully")
	return &MinioStorageService{
		config: config,
		client: client,
	}, nil
}

func (s *MinioStorageService) DownloadFile(ctx context.Context, url string) (string, error) {
	logger.Debug("downloading file", "url", url)

	// Create temporary file
	tmpFile, err := os.CreateTemp(s.config.Storage.TempDirectory, "input-*")
	if err != nil {
		logger.Error("failed to create temp file", "error", err)
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// For external URLs, download using HTTP
	if isExternalURL(url) {
		logger.Debug("downloading from external URL", "url", url)
		resp, err := http.Get(url)
		if err != nil {
			os.Remove(tmpFile.Name())
			logger.Error("failed to download from external URL", "error", err)
			return "", fmt.Errorf("failed to download file: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			os.Remove(tmpFile.Name())
			logger.Error("failed to download file",
				"status_code", resp.StatusCode,
				"url", url)
			return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
		}

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			os.Remove(tmpFile.Name())
			logger.Error("failed to save downloaded file", "error", err)
			return "", fmt.Errorf("failed to save file: %w", err)
		}
	} else {
		// For MinIO objects, get the object
		logger.Debug("downloading from MinIO",
			"bucket", s.config.Storage.MinioBucketName,
			"object", url)
		object, err := s.client.GetObject(ctx, s.config.Storage.MinioBucketName, url, minio.GetObjectOptions{})
		if err != nil {
			os.Remove(tmpFile.Name())
			logger.Error("failed to get object from MinIO", "error", err)
			return "", fmt.Errorf("failed to get object from MinIO: %w", err)
		}
		defer object.Close()

		if _, err := io.Copy(tmpFile, object); err != nil {
			os.Remove(tmpFile.Name())
			logger.Error("failed to save MinIO object", "error", err)
			return "", fmt.Errorf("failed to save file: %w", err)
		}
	}

	logger.Debug("file downloaded successfully", "path", tmpFile.Name())
	return tmpFile.Name(), nil
}

func (s *MinioStorageService) UploadFile(ctx context.Context, localPath string, objectKey string, userID uint) (string, error) {
	logger.Debug("uploading file",
		"local_path", localPath,
		"object_key", objectKey,
		"user_id", userID)

	// Open the local file
	file, err := os.Open(localPath)
	if err != nil {
		logger.Error("failed to open file", "error", err)
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for content-type detection
	fileInfo, err := file.Stat()
	if err != nil {
		logger.Error("failed to get file info", "error", err)
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Create user-specific object key
	userObjectKey := fmt.Sprintf("user_%d/%s", userID, objectKey)

	// Upload the file to MinIO
	logger.Debug("uploading to MinIO",
		"bucket", s.config.Storage.MinioBucketName,
		"object", userObjectKey,
		"size", fileInfo.Size())

	_, err = s.client.PutObject(ctx, s.config.Storage.MinioBucketName, userObjectKey, file, fileInfo.Size(),
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			UserMetadata: map[string]string{
				"x-amz-meta-filename": objectKey,
				"x-amz-meta-userid":   fmt.Sprintf("%d", userID),
			},
			Expires: time.Now().Add(time.Hour),
		})
	if err != nil {
		logger.Error("failed to upload file to MinIO", "error", err)
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	logger.Info("file uploaded successfully",
		"bucket", s.config.Storage.MinioBucketName,
		"object", userObjectKey)

	// return the full url
	return fmt.Sprintf("%s/%s", s.config.Storage.MinioBucketURL, userObjectKey), nil
}

func (s *MinioStorageService) DeleteFile(ctx context.Context, localPath string) error {
	logger.Debug("deleting file", "path", localPath)
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		logger.Error("failed to delete file", "error", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// isExternalURL checks if the URL is external (starts with http:// or https://)
func isExternalURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}
