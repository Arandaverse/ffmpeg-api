package repository

import (
	"ffmpeg-api/internal/database"

	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db database.Database
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db database.Database) BaseRepository {
	return BaseRepository{db: db}
}

// GetDB returns the underlying database instance
func (r *BaseRepository) GetDB() database.Database {
	return r.db
}

// GetGormDB returns the underlying GORM database instance
func (r *BaseRepository) GetGormDB() (*gorm.DB, error) {
	if gormDB, ok := r.db.DB().(*gorm.DB); ok {
		return gormDB, nil
	}
	return nil, database.ErrNotGormDB
}
