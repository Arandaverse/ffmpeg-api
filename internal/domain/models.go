package domain

import "time"

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

// JobStatus represents the status of an FFMPEG job.
type JobStatus struct {
	ID                      uint                          `gorm:"primaryKey" json:"id"`
	UUID                    string                        `gorm:"uniqueIndex" json:"uuid"`
	Status                  string                        `json:"status"`
	Result                  string                        `json:"-"`
	Progress                int                           `json:"progress"`
	Error                   string                        `json:"error,omitempty"`
	UserID                  uint                          `json:"user_id"`
	OriginalRequest         *FFMPEGRequest                `json:"original_request,omitempty" gorm:"-"`
	OutputFiles             map[string]OutputFileMetadata `json:"output_files,omitempty" gorm:"-"`
	FFmpegCommandRunSeconds float64                       `json:"ffmpeg_command_run_seconds,omitempty"`
	TotalProcessingSeconds  float64                       `json:"total_processing_seconds,omitempty"`
	CreatedAt               time.Time                     `json:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at"`
}

// FFMPEGRequest represents the request body for the /ffmpeg endpoint.
type FFMPEGRequest struct {
	InputFiles    map[string]string `json:"input_files"`
	OutputFiles   map[string]string `json:"output_files"`
	FFmpegCommand string            `json:"ffmpeg_command"`
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
