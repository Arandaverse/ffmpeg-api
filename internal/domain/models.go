package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// User represents a user in the system.
type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"uniqueIndex" json:"username"`
	Email          string    `gorm:"uniqueIndex" json:"email"`
	Password       string    `json:"-"`
	APIToken       string    `gorm:"uniqueIndex" json:"api_token"`
	UsageCount     int       `gorm:"default:0" json:"usage_count"`
	BytesProcessed int64     `gorm:"default:0" json:"bytes_processed"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// OutputFileMetadata represents metadata for a processed output file
type OutputFileMetadata struct {
	FileID     string  `json:"file_id"`
	SizeMBytes float64 `json:"size_mbytes"`
	FileType   string  `json:"file_type"`
	FileFormat string  `json:"file_format"`
	StorageURL string  `json:"storage_url"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
}

// OutputFilesMap is a custom type for handling the map of output files in the database
type OutputFilesMap map[string]OutputFileMetadata

// Scan implements the sql.Scanner interface for OutputFilesMap
func (o *OutputFilesMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, &o)
}

// Value implements the driver.Valuer interface for OutputFilesMap
func (o OutputFilesMap) Value() (driver.Value, error) {
	if o == nil {
		return nil, nil
	}
	return json.Marshal(o)
}

func (o *OutputFileMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), &o)
}

// JobStatus represents the status of an FFMPEG job.
type JobStatus struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	UUID                    string         `gorm:"uniqueIndex" json:"uuid"`
	Status                  string         `json:"status"`
	Result                  string         `json:"-"`
	Progress                int            `json:"progress"`
	Error                   string         `json:"error,omitempty"`
	UserID                  uint           `json:"user_id"`
	OriginalRequest         *FFMPEGRequest `json:"original_request,omitempty" gorm:"type:jsonb"`
	OutputFiles             OutputFilesMap `json:"output_files,omitempty" gorm:"type:jsonb"`
	FFmpegCommandRunSeconds float64        `json:"ffmpeg_command_run_seconds,omitempty"`
	TotalProcessingSeconds  float64        `json:"total_processing_seconds,omitempty"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}

// FFMPEGRequest represents the request body for the /ffmpeg endpoint.
type FFMPEGRequest struct {
	InputFiles    map[string]string `json:"input_files" gorm:"type:jsonb"`
	OutputFiles   map[string]string `json:"output_files" gorm:"type:jsonb"`
	FFmpegCommand string            `json:"ffmpeg_command"`
}

// Scan implements the sql.Scanner interface for FFMPEGRequest
func (f *FFMPEGRequest) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	if err := json.Unmarshal(bytes, &f); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for FFMPEGRequest
func (f FFMPEGRequest) Value() (driver.Value, error) {
	if f.InputFiles == nil && f.OutputFiles == nil && f.FFmpegCommand == "" {
		return nil, nil
	}
	return json.Marshal(f)
}

// FFMPEGResponse represents the response from the FFMPEG processing endpoint.
type FFMPEGResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

// RegisterRequest represents the user registration request.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginRequest represents the user login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response.
type AuthResponse struct {
	APIToken string `json:"api_token"`
}
