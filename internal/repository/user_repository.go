package repository

import (
	"context"
	"ffmpeg-api/internal/database"
	"ffmpeg-api/internal/domain"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	BaseRepository
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db database.Database) UserRepository {
	return &GormUserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(user).Error
}

func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *GormUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := db.WithContext(ctx).Where("username = ?", username).Select("username, email, password, api_token").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByAPIToken(ctx context.Context, token string) (*domain.User, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := db.WithContext(ctx).Where("api_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) IncrementUsage(ctx context.Context, userID uint) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *GormUserRepository) IncrementBytesProcessed(ctx context.Context, userID uint, bytes int64) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("bytes_processed", gorm.Expr("bytes_processed + ?", bytes)).Error
}
