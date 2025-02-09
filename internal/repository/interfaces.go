package repository

import (
	"context"
	"ffmpeg-api/internal/domain"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error)
	FindByAPIToken(ctx context.Context, token string) (*domain.User, error)
	IncrementUsage(ctx context.Context, userID uint) error
	IncrementBytesProcessed(ctx context.Context, userID uint, bytes int64) error
}

// JobRepository defines the interface for job-related database operations
type JobRepository interface {
	Create(ctx context.Context, job *domain.JobStatus) error
	Update(ctx context.Context, job *domain.JobStatus) error
	FindByUUID(ctx context.Context, uuid string) (*domain.JobStatus, error)
	FindByUserID(ctx context.Context, userID uint) ([]domain.JobStatus, error)
}
