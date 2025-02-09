package repository

import (
	"context"
	"ffmpeg-api/internal/database"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/logger"
)

type GormJobRepository struct {
	BaseRepository
}

func NewGormJobRepository(db database.Database) JobRepository {
	return &GormJobRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *GormJobRepository) Create(ctx context.Context, job *domain.JobStatus) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(job).Error
}

func (r *GormJobRepository) FindByID(ctx context.Context, id uint) (*domain.JobStatus, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var job domain.JobStatus
	if err := db.WithContext(ctx).First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *GormJobRepository) FindByUUID(ctx context.Context, uuid string) (*domain.JobStatus, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	logger.Debug("Finding job by UUID: " + uuid)
	var job domain.JobStatus
	if err := db.WithContext(ctx).Where("uuid = ?", uuid).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *GormJobRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.JobStatus, error) {
	db, err := r.GetGormDB()
	if err != nil {
		return nil, err
	}
	var jobs []domain.JobStatus
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *GormJobRepository) Update(ctx context.Context, job *domain.JobStatus) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(job).Error
}

func (r *GormJobRepository) Delete(ctx context.Context, id uint) error {
	db, err := r.GetGormDB()
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&domain.JobStatus{}, id).Error
}
