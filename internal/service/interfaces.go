package service

import (
	"context"
	"ffmpeg-api/internal/domain"
)

// AuthService defines the interface for authentication-related operations
type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

// FFMPEGService defines the interface for FFMPEG processing operations
type FFMPEGService interface {
	ProcessVideo(ctx context.Context, req domain.FFMPEGRequest, userID uint) (*domain.FFMPEGResponse, error)
	GetJobStatus(ctx context.Context, uuid string, userID uint) (*domain.JobStatus, error)
}

// StorageService defines the interface for file storage operations
type StorageService interface {
	DownloadFile(ctx context.Context, url string) (string, error)
	UploadFile(ctx context.Context, localPath string, objectKey string, userID uint) (string, error)
	DeleteFile(ctx context.Context, localPath string) error
}
