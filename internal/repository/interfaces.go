package repository

import (
	"context"
	"ffmpeg-api/internal/domain"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	BaseRepositoryInterface[domain.User]
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error)
	FindByAPIToken(ctx context.Context, token string) (*domain.User, error)
	IncrementUsage(ctx context.Context, userID uint) error
	IncrementBytesProcessed(ctx context.Context, userID uint, bytes int64) error
}

type JobRepository interface {
	BaseRepositoryInterface[domain.JobStatus]
	FindByUUID(ctx context.Context, uuid string) (*domain.JobStatus, error)
	FindByUserID(ctx context.Context, userID uint) ([]domain.JobStatus, error)
}
