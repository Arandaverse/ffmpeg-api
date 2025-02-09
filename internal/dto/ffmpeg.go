package dto

import "ffmpeg-api/internal/domain"

// FFMPEGRequest represents the FFMPEG processing request
type FFMPEGRequest struct {
	InputFiles    map[string]string `json:"input_files" validate:"required,min=1" example:"{\"in1\": \"https://storage.googleapis.com/ffmpeg-api-test-bucket/user_1/input/test.mp4\"}"`
	OutputFiles   map[string]string `json:"output_files" validate:"required,min=1" example:"{\"out1\": \"string.mp4\"}"`
	FFmpegCommand string            `json:"ffmpeg_command" validate:"required" example:"-i {{in1}} {{out1}}"`
}

// FFMPEGResponse represents the FFMPEG processing response
type FFMPEGResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

// JobStatus represents the status of an FFMPEG job
type JobStatus struct {
	UUID        string                               `json:"uuid"`
	Status      string                               `json:"status" validate:"required,oneof=pending processing completed failed"`
	Result      string                               `json:"result,omitempty"`
	Progress    int                                  `json:"progress"`
	Error       string                               `json:"error,omitempty"`
	CreatedAt   string                               `json:"created_at"`
	UpdatedAt   string                               `json:"updated_at"`
	OutputFiles map[string]domain.OutputFileMetadata `json:"output_files,omitempty"`
}
