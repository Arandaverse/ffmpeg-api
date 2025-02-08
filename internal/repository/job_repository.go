package repository

import (
	"context"
	"ffmpeg-api/internal/domain"

	"gorm.io/gorm"
)

// GormJobRepository implements JobRepository using GORM
type GormJobRepository struct {
	db *gorm.DB
}

// NewGormJobRepository creates a new GormJobRepository
func NewGormJobRepository(db *gorm.DB) *GormJobRepository {
	return &GormJobRepository{db: db}
}

func (r *GormJobRepository) Create(ctx context.Context, job *domain.JobStatus) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *GormJobRepository) Update(ctx context.Context, job *domain.JobStatus) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *GormJobRepository) FindByUUID(ctx context.Context, uuid string) (*domain.JobStatus, error) {
	var job domain.JobStatus
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *GormJobRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.JobStatus, error) {
	var jobs []domain.JobStatus
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
