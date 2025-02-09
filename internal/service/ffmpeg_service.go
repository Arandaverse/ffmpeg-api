package service

import (
	"context"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/repository"
	"fmt"
	"image"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"bufio"

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
	jobUUID := uuid.New().String()

	job := &domain.JobStatus{
		UUID:   jobUUID,
		Status: "pending",
		UserID: userID,
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

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

	if job.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to job")
	}

	return job, nil
}

func (s *FFMPEGServiceImpl) processFFMPEGJob(ctx context.Context, job *domain.JobStatus, req domain.FFMPEGRequest) {
	startTime := time.Now()
	job.Status = "PROCESSING"
	job.Progress = 0 // Initialize progress
	job.OriginalRequest = &req
	if err := s.jobRepo.Update(ctx, job); err != nil {
		s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to update job status: %v", err))
		return
	}

	// Create temporary directory for this job
	tempDir := filepath.Join(s.config.FFMPEG.TempDirectory, job.UUID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to create temp directory: %v", err))
		return
	}
	defer os.RemoveAll(tempDir)

	// Download all input files (25% of progress)
	inputPaths := make(map[string]string)
	var totalInputSize int64
	totalFiles := len(req.InputFiles)
	fileNum := 0
	for key, url := range req.InputFiles {
		if !strings.HasPrefix(url, "http") {
			s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("Invalid URL for input file %s", key))
			return
		}

		inputPath, err := s.storageService.DownloadFile(ctx, url)
		if err != nil {
			s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to download input file %s: %v", key, err))
			return
		}
		defer s.storageService.DeleteFile(ctx, inputPath)

		inputFileInfo, err := os.Stat(inputPath)
		if err != nil {
			s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to get input file size for %s: %v", key, err))
			return
		}
		totalInputSize += inputFileInfo.Size()
		inputPaths[key] = inputPath

		// Update progress for download phase (0-25%)
		fileNum++
		job.Progress = int(float64(fileNum) / float64(totalFiles) * 25)
		if err := s.jobRepo.Update(ctx, job); err != nil {
			logger.Error("failed to update job progress", "error", err)
		}
	}

	// Prepare output paths
	outputPaths := make(map[string]string)
	for key, filename := range req.OutputFiles {
		outputPaths[key] = filepath.Join(tempDir, filename)
	}

	// Process command template
	command := req.FFmpegCommand
	for key, path := range inputPaths {
		placeholder := fmt.Sprintf("{{%s}}", key)
		logger.Debug(command, "Command<<")
		logger.Debug(placeholder, "Placeholder<<")
		logger.Debug(path, "Path<<")
		command = strings.ReplaceAll(command, placeholder, path)
		logger.Debug(command, "Command2<<")
		logger.Debug(path, "Path2<<")
	}
	for key, path := range outputPaths {
		placeholder := fmt.Sprintf("{{%s}}", key)
		command = strings.ReplaceAll(command, placeholder, path)
		logger.Debug(command, "Command3<<")
		// println(path, "Path<<")
	}

	// Split command into args
	args := splitCommand(command)
	fmt.Println(args, "<<<<")

	if len(args) == 0 {
		s.updateJobStatus(ctx, job, "FAILED", "invalid FFmpeg command")
		return
	}

	// Execute FFmpeg command (25-75% of progress)
	job.Progress = 25
	if err := s.jobRepo.Update(ctx, job); err != nil {
		logger.Error("failed to update job progress", "error", err)
	}

	ffmpegStartTime := time.Now()
	cmd := exec.CommandContext(ctx, s.config.FFMPEG.BinaryPath, args...)

	// Capture stderr to parse progress
	stderr, err := cmd.StderrPipe()
	if err != nil {
		s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to create stderr pipe: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to start FFmpeg: %v", err))
		return
	}

	// Start a goroutine to read stderr and update progress
	service := s // Capture service instance for goroutine
	go func() {
		scanner := bufio.NewScanner(stderr)
		var duration float64

		// First, try to find the duration
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Duration:") {
				// Parse duration in format "Duration: 00:00:00.00"
				parts := strings.Split(line, "Duration: ")
				if len(parts) > 1 {
					timeStr := strings.Split(parts[1], ",")[0]
					h, m, s := 0, 0, 0.0
					fmt.Sscanf(timeStr, "%d:%d:%f", &h, &m, &s)
					duration = float64(h*3600) + float64(m*60) + s
					break
				}
			}
		}

		// Now process the time updates
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "time=") {
				// Parse time in format "time=00:00:00.00"
				parts := strings.Split(line, "time=")
				if len(parts) > 1 {
					timeStr := strings.Split(parts[1], " ")[0]
					h, m, s := 0, 0, 0.0
					fmt.Sscanf(timeStr, "%d:%d:%f", &h, &m, &s)
					currentTime := float64(h*3600) + float64(m*60) + s

					if duration > 0 {
						// Calculate progress within the FFMPEG phase (25-75%)
						ffmpegProgress := (currentTime / duration) * 50
						job.Progress = 25 + int(ffmpegProgress)
						if err := service.jobRepo.Update(ctx, job); err != nil {
							logger.Error("failed to update job progress", "error", err)
						}
					}
				}
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("FFmpeg processing failed: %v", err))
		return
	}
	ffmpegEndTime := time.Now()
	job.FFmpegCommandRunSeconds = ffmpegEndTime.Sub(ffmpegStartTime).Seconds()

	// Upload output files and gather metadata (75-99% of progress)
	job.Progress = 75
	if err := s.jobRepo.Update(ctx, job); err != nil {
		logger.Error("failed to update job progress", "error", err)
	}

	var totalOutputSize int64
	job.OutputFiles = make(map[string]domain.OutputFileMetadata)

	totalOutputFiles := len(outputPaths)
	currentOutputFile := 0

	for key, outputPath := range outputPaths {
		outputFileInfo, err := os.Stat(outputPath)
		if err != nil {
			s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to get output file size for %s: %v", key, err))
			return
		}
		totalOutputSize += outputFileInfo.Size()

		s3URL, err := s.storageService.UploadFile(ctx, outputPath, filepath.Base(outputPath), job.UserID)
		if err != nil {
			s.updateJobStatus(ctx, job, "FAILED", fmt.Sprintf("failed to upload output file %s: %v", key, err))
			return
		}

		logger.Debug(s3URL, "S3 URL<<")

		// Get file metadata
		metadata := domain.OutputFileMetadata{
			FileID:     uuid.New().String(),
			SizeMBytes: float64(outputFileInfo.Size()) / 1024 / 1024,
			StorageURL: s3URL,
		}

		// Get file format from extension
		ext := strings.TrimPrefix(filepath.Ext(outputPath), ".")
		metadata.FileFormat = ext

		// Set file type based on format
		switch ext {
		case "jpg", "jpeg", "png", "gif", "webp":
			metadata.FileType = "image"
			// Get image dimensions
			if f, err := os.Open(outputPath); err == nil {
				defer f.Close()
				if img, _, err := image.DecodeConfig(f); err == nil {
					metadata.Width = img.Width
					metadata.Height = img.Height
				}
			}
		case "mp4", "webm", "mov", "avi":
			metadata.FileType = "video"
			// For video files, we could use ffprobe here if needed
		default:
			metadata.FileType = "unknown"
		}

		job.OutputFiles[key] = metadata

		currentOutputFile++
		// Update progress for upload phase (75-99%)
		job.Progress = 75 + int(float64(currentOutputFile)/float64(totalOutputFiles)*24)
		// add files to job

		if err := s.jobRepo.Update(ctx, job); err != nil {
			logger.Error("failed to update job progress", "error", err)
		}
	}

	// Update job status to completed and set progress to 100%
	job.Status = "SUCCESS"
	job.Progress = 100
	job.TotalProcessingSeconds = time.Since(startTime).Seconds()
	job.Result = "Successfully processed files"
	if err := s.jobRepo.Update(ctx, job); err != nil {
		fmt.Printf("failed to update final job status: %v\n", err)
		return
	}

	// Update user usage statistics
	go s.userRepo.IncrementUsage(ctx, job.UserID)
	go s.userRepo.IncrementBytesProcessed(ctx, job.UserID, totalInputSize+totalOutputSize)
}

func (s *FFMPEGServiceImpl) updateJobStatus(ctx context.Context, job *domain.JobStatus, status, result string) {
	job.Status = status
	job.Result = result
	job.UpdatedAt = time.Now()
	job.TotalProcessingSeconds = time.Since(job.CreatedAt).Seconds()

	if err := s.jobRepo.Update(ctx, job); err != nil {
		fmt.Printf("failed to update job status: %v\n", err)
	}
}

// splitCommand splits a command string into arguments, respecting quotes
func splitCommand(command string) []string {
	r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)
	matches := r.FindAllString(command, -1)
	args := make([]string, 0, len(matches))
	for _, match := range matches {
		// Remove surrounding quotes if present
		if (strings.HasPrefix(match, "\"") && strings.HasSuffix(match, "\"")) ||
			(strings.HasPrefix(match, "'") && strings.HasSuffix(match, "'")) {
			match = match[1 : len(match)-1]
		}
		args = append(args, match)
	}
	return args
}
