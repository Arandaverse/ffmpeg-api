package dto

// FFMPEGRequest represents the FFMPEG processing request
type FFMPEGRequest struct {
	Command   string `json:"command" validate:"required"`
	S3FileURL string `json:"s3_file_url" validate:"required,url"`
	Format    string `json:"format" validate:"required,oneof=mp4 webm mov avi"`
	Quality   string `json:"quality" validate:"required,oneof=low medium high"`
}

// FFMPEGResponse represents the FFMPEG processing response
type FFMPEGResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

// JobStatus represents the status of an FFMPEG job
type JobStatus struct {
	UUID      string `json:"uuid"`
	Status    string `json:"status" validate:"required,oneof=pending processing completed failed"`
	Result    string `json:"result,omitempty"`
	Progress  int    `json:"progress"`
	Error     string `json:"error,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
