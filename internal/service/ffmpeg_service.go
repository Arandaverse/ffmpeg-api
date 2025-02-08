package service

import (
	"context"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/repository"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// FFMPEGServiceImpl implements FFMPEGService
type FFMPEGServiceImpl struct {
	jobRepo        repository.JobRepository
	userRepo       repository.UserRepository
	storageService StorageService
	config         *config.Config
}

// NewFFMPEGService creates a new FFMPEGService
func NewFFMPEGService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	storageService StorageService,
	config *config.Config,
) FFMPEGService {
	return &FFMPEGServiceImpl{
		jobRepo:        jobRepo,
		userRepo:       userRepo,
		storageService: storageService,
		config:         config,
	}
}

func (s *FFMPEGServiceImpl) ProcessVideo(ctx context.Context, req domain.FFMPEGRequest, userID uint) (*domain.FFMPEGResponse, error) {
	// Generate UUID for the job
	jobUUID := uuid.New().String()

	// Create job status
	job := &domain.JobStatus{
		UUID:   jobUUID,
		Status: "pending",
		UserID: userID,
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Start processing in background
	go s.processFFMPEGJob(context.Background(), job, req)

	return &domain.FFMPEGResponse{
		UUID:   jobUUID,
		Status: "pending",
	}, nil
}

func (s *FFMPEGServiceImpl) GetJobStatus(ctx context.Context, uuid string, userID uint) (*domain.JobStatus, error) {
	job, err := s.jobRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Verify the job belongs to the user
	if job.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to job")
	}

	return job, nil
}

func (s *FFMPEGServiceImpl) processFFMPEGJob(ctx context.Context, job *domain.JobStatus, req domain.FFMPEGRequest) {
	// Update job status to processing
	job.Status = "processing"
	if err := s.jobRepo.Update(ctx, job); err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to update job status: %v", err))
		return
	}

	// Download input file
	inputPath, err := s.storageService.DownloadFile(ctx, req.S3FileURL)
	if err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to download input file: %v", err))
		return
	}
	defer s.storageService.DeleteFile(ctx, inputPath)

	// Get input file size
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to get input file size: %v", err))
		return
	}
	inputSize := inputFileInfo.Size()

	// Prepare output file path
	outputFileName := fmt.Sprintf("%s%s", job.UUID, filepath.Ext(inputPath))
	outputPath := filepath.Join(s.config.FFMPEG.TempDirectory, outputFileName)

	// Ensure temp directory exists
	if err := os.MkdirAll(s.config.FFMPEG.TempDirectory, 0755); err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to create temp directory: %v", err))
		return
	}

	// Prepare FFMPEG command
	cmd := exec.CommandContext(ctx, s.config.FFMPEG.BinaryPath, "-i", inputPath, "-y", outputPath)

	// Start command
	if err := cmd.Start(); err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to start FFMPEG: %v", err))
		return
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("FFMPEG processing failed: %v", err.Error()))
		return
	}

	// Get output file size
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to get output file size: %v", err))
		return
	}
	outputSize := outputFileInfo.Size()

	// Upload processed file
	s3URL, err := s.storageService.UploadFile(ctx, outputPath, outputFileName)
	if err != nil {
		s.updateJobStatus(ctx, job, "failed", fmt.Sprintf("failed to upload processed file: %v", err))
		return
	}

	// Clean up output file
	// defer s.storageService.DeleteFile(ctx, outputPath)

	// Update job status to completed
	s.updateJobStatus(ctx, job, "completed", s3URL)

	// Increment user usage count and bytes processed
	go s.userRepo.IncrementUsage(ctx, job.UserID)

	// Track total bytes processed (input + output)
	totalBytes := inputSize + outputSize
	go s.userRepo.IncrementBytesProcessed(ctx, job.UserID, totalBytes)
}

func (s *FFMPEGServiceImpl) updateJobStatus(ctx context.Context, job *domain.JobStatus, status, result string) {
	job.Status = status
	job.Result = result
	job.UpdatedAt = time.Now()

	if err := s.jobRepo.Update(ctx, job); err != nil {
		fmt.Printf("failed to update job status: %v\n", err)
	}
}
