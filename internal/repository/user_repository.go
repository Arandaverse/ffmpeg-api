package repository

import (
	"context"
	"ffmpeg-api/internal/domain"

	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).First(&user).Error
}

func (r *GormUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).Select("username, email, password, api_token").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByAPIToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("api_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) IncrementUsage(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *GormUserRepository) IncrementBytesProcessed(ctx context.Context, userID uint, bytes int64) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("bytes_processed", gorm.Expr("bytes_processed + ?", bytes)).Error
}
