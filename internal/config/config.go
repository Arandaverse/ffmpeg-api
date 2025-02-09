package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	FFMPEG   FFMPEGConfig
	Storage  StorageConfig
}

// ServerConfig holds HTTP server related configuration
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	APITokenLength int
	AllowedOrigins []string
}

// DatabaseConfig holds database related configuration
type DatabaseConfig struct {
	Driver string
	DSN    string
}

// FFMPEGConfig holds FFMPEG related configuration
type FFMPEGConfig struct {
	BinaryPath             string
	TempDirectory          string
	ProgressUpdateInterval time.Duration
}

// StorageConfig holds storage related configuration
type StorageConfig struct {
	Provider        string // "local" or "minio"
	TempDirectory   string
	MinioEndpoint   string
	MinioPort       string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
	MinioBucketName string
	MinioRegion     string
	MinioBucketURL  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	apiTokenLength, _ := strconv.Atoi(getEnv("API_TOKEN_LENGTH", "32"))
	progressInterval, _ := strconv.Atoi(getEnv("PROGRESS_UPDATE_INTERVAL", "5"))
	useSSL, _ := strconv.ParseBool(getEnv("MINIO_USE_SSL", "false"))

	return &Config{
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8000"),
			ReadTimeout:    time.Second * 15,
			WriteTimeout:   time.Second * 15,
			APITokenLength: apiTokenLength,
			AllowedOrigins: []string{"*"}, // Configure as needed
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "ffmpeg_api.db"),
		},
		FFMPEG: FFMPEGConfig{
			BinaryPath:             getEnv("FFMPEG_PATH", "/usr/bin/ffmpeg"),
			TempDirectory:          getEnv("TEMP_DIR", "tmp"),
			ProgressUpdateInterval: time.Duration(progressInterval) * time.Second,
		},
		Storage: StorageConfig{
			Provider:        getEnv("STORAGE_PROVIDER", "local"),
			TempDirectory:   getEnv("TEMP_DIR", "tmp"),
			MinioEndpoint:   getEnv("MINIO_ENDPOINT", "127.0.0.1"),
			MinioPort:       getEnv("MINIO_PORT", "56732"),
			MinioAccessKey:  getEnv("MINIO_ACCESS_KEY", ""),
			MinioSecretKey:  getEnv("MINIO_SECRET_KEY", ""),
			MinioUseSSL:     useSSL,
			MinioBucketName: getEnv("MINIO_BUCKET_NAME", "ffmpeg-files"),
			MinioRegion:     getEnv("MINIO_REGION", "us-east-1"),
			MinioBucketURL:  getEnv("MINIO_BUCKET_URL", "http://127.0.0.1:9000"),
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
