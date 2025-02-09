package routes

import (
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/dto"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/response"
	"ffmpeg-api/internal/service"
	"ffmpeg-api/internal/validation"

	"github.com/gofiber/fiber/v2"
)

// FFMPEGRoutes handles all FFMPEG related routes
type FFMPEGRoutes struct {
	ffmpegService service.FFMPEGService
	authService   service.AuthService
}

// NewFFMPEGRoutes creates a new FFMPEGRoutes instance
func NewFFMPEGRoutes(ffmpegService service.FFMPEGService, authService service.AuthService) *FFMPEGRoutes {
	return &FFMPEGRoutes{
		ffmpegService: ffmpegService,
		authService:   authService,
	}
}

// Register registers all FFMPEG routes
func (r *FFMPEGRoutes) Register(router fiber.Router) {
	ffmpeg := router.Group("/ffmpeg")
	ffmpeg.Use(r.authMiddleware)
	ffmpeg.Post("/", r.handleProcessFFMPEG)
	ffmpeg.Get("/progress/:uuid", r.handleGetProgress)
}

// handleProcessFFMPEG handles video processing requests
// @Summary Process video with FFMPEG
// @Description Submit a video processing job using FFMPEG
// @Tags FFMPEG
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.FFMPEGRequest true "FFMPEG processing details"
// @Success 202 {object} response.Response{data=dto.FFMPEGResponse}
// @Failure 400 {object} response.Response{error=response.APIError}
// @Failure 401 {object} response.Response{error=response.APIError}
// @Failure 500 {object} response.Response{error=response.APIError}
// @Router /ffmpeg [post]
func (r *FFMPEGRoutes) handleProcessFFMPEG(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		logger.Error("user not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Unauthorized",
				Message: "Unauthorized",
			},
		})
	}

	var req dto.FFMPEGRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("invalid request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "BadRequest",
				Message: "Invalid request body",
			},
		})
	}

	if err := validation.Validate(req); err != nil {
		logger.Error("validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "ValidationError",
				Message: err.Error(),
			},
		})
	}

	// Convert DTO to domain model
	domainReq := domain.FFMPEGRequest{
		InputFiles:    req.InputFiles,
		OutputFiles:   req.OutputFiles,
		FFmpegCommand: req.FFmpegCommand,
	}

	resp, err := r.ffmpegService.ProcessVideo(c.Context(), domainReq, user.ID)
	if err != nil {
		logger.Error("failed to process video", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "InternalServerError",
				Message: "Failed to process video",
			},
		})
	}

	// Convert domain response to DTO
	dtoResp := dto.FFMPEGResponse{
		UUID:   resp.UUID,
		Status: resp.Status,
	}

	return c.Status(fiber.StatusAccepted).JSON(response.Response{
		Success: true,
		Data:    dtoResp,
	})
}

// handleGetProgress handles job progress requests
// @Summary Get job progress
// @Description Get the progress of a video processing job
// @Tags FFMPEG
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "Job UUID"
// @Success 200 {object} response.Response{data=dto.JobStatus}
// @Failure 400 {object} response.Response{error=response.APIError}
// @Failure 401 {object} response.Response{error=response.APIError}
// @Failure 404 {object} response.Response{error=response.APIError}
// @Router /ffmpeg/progress/{uuid} [get]
func (r *FFMPEGRoutes) handleGetProgress(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		logger.Error("user not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Unauthorized",
				Message: "Unauthorized",
			},
		})
	}

	uuid := c.Params("uuid")
	if uuid == "" {
		logger.Error("missing uuid parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "BadRequest",
				Message: "Missing UUID parameter",
			},
		})
	}

	status, err := r.ffmpegService.GetJobStatus(c.Context(), uuid, user.ID)
	if err != nil {
		logger.Error("failed to get job status", "error", err, "uuid", uuid)
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "NotFound",
				Message: "Job not found",
			},
		})
	}

	// Convert domain model to DTO
	dtoStatus := dto.JobStatus{
		UUID:        status.UUID,
		Status:      status.Status,
		Result:      status.Result,
		Progress:    status.Progress,
		Error:       status.Error,
		CreatedAt:   status.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   status.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		OutputFiles: status.OutputFiles,
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Success: true,
		Data:    dtoStatus,
	})
}

// authMiddleware authenticates requests
func (r *FFMPEGRoutes) authMiddleware(c *fiber.Ctx) error {
	token := c.Get("X-API-Token")
	if token == "" {
		logger.Warn("missing API token")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Unauthorized",
				Message: "Missing API token",
			},
		})
	}

	user, err := r.authService.ValidateToken(c.Context(), token)
	if err != nil {
		logger.Error("invalid API token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Unauthorized",
				Message: "Invalid API token",
			},
		})
	}

	c.Locals("user", user)
	return c.Next()
}
