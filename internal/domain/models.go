package domain

import "time"

// User represents a user in the system.
type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"uniqueIndex" json:"username"`
	Email          string    `gorm:"uniqueIndex" json:"email"`
	APIToken       string    `gorm:"uniqueIndex" json:"api_token"`
	UsageCount     int       `gorm:"default:0" json:"usage_count"`
	BytesProcessed int64     `gorm:"default:0" json:"bytes_processed"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// JobStatus represents the status of an FFMPEG job.
type JobStatus struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"uniqueIndex" json:"uuid"`
	Status    string    `json:"status"`
	Result    string    `json:"result"`
	Progress  int       `json:"progress"`
	Error     string    `json:"error"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FFMPEGRequest represents the request body for the /ffmpeg endpoint.
type FFMPEGRequest struct {
	Command   string `json:"command"`
	S3FileURL string `json:"s3_file_url"`
	Format    string `json:"format"`
	Quality   string `json:"quality"`
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
