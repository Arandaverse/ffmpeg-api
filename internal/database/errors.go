package database

import "errors"

var (
	// ErrNotGormDB is returned when trying to get a GORM DB instance from a non-GORM database
	ErrNotGormDB = errors.New("database is not a GORM database")
)
