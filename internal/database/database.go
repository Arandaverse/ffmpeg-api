package database

import (
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/logger"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database represents the interface for database operations
type Database interface {
	// DB returns the underlying database instance
	DB() interface{}
	// AutoMigrate runs database migrations for the given models
	AutoMigrate(models ...interface{}) error
	// Close closes the database connection
	Close() error
}

// GormDatabase implements the Database interface using GORM
type GormDatabase struct {
	db *gorm.DB
}

// NewDatabase creates a new database instance based on the configuration
func NewDatabase(cfg *config.Config) (Database, error) {
	switch cfg.Database.Driver {
	case "sqlite":
		return newSQLiteDatabase(cfg)
	// Add more database drivers here
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
}

func newSQLiteDatabase(cfg *config.Config) (Database, error) {
	logger.Info("initializing SQLite database")
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	return &GormDatabase{db: db}, nil
}

// DB returns the underlying GORM database instance
func (d *GormDatabase) DB() interface{} {
	return d.db
}

// AutoMigrate runs database migrations for the given models
func (d *GormDatabase) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

// Close closes the database connection
func (d *GormDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetGormDB safely converts the Database interface to a *gorm.DB instance
func GetGormDB(db Database) (*gorm.DB, error) {
	if gormDB, ok := db.DB().(*gorm.DB); ok {
		return gormDB, nil
	}
	return nil, fmt.Errorf("database is not a GORM database")
}
