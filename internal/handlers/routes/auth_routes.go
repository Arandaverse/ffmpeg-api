package routes

import (
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/dto"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/response"
	"ffmpeg-api/internal/service"
	"ffmpeg-api/internal/validation"
	"os"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes handles all authentication related routes
type AuthRoutes struct {
	authService service.AuthService
}

// NewAuthRoutes creates a new AuthRoutes instance
func NewAuthRoutes(authService service.AuthService) *AuthRoutes {
	return &AuthRoutes{
		authService: authService,
	}
}

// Register registers all auth routes
func (r *AuthRoutes) Register(router fiber.Router) {
	auth := router.Group("/api/v1/auth")
	auth.Post("/register", r.handleRegister)
	auth.Post("/login", r.handleLogin)
}

// handleRegister handles user registration
// @Summary Register a new user
// @Description Register a new user account with username, password and email. The password must be at least 8 characters long.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=dto.AuthResponse} "Successfully registered"
// @Failure 400 {object} response.Response{error=response.APIError} "Invalid request or validation error"
// @Failure 409 {object} response.Response{error=response.APIError} "Username or email already exists"
// @Failure 500 {object} response.Response{error=response.APIError} "Internal server error"
// @Router /auth/register [post]
func (r *AuthRoutes) handleRegister(c *fiber.Ctx) error {
	var req dto.RegisterRequest
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

	if req.RegisterKey != os.Getenv("REGISTER_KEY") {
		logger.Error("Invalid register key", "error")
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "BadRequest",
				Message: "Invalid register key",
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
	domainReq := domain.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	resp, err := r.authService.Register(c.Context(), domainReq)
	if err != nil {
		logger.Error(err.Error(), "error", err)
		return c.Status(fiber.StatusConflict).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Conflict",
				Message: err.Error(),
			},
		})
	}

	// Convert domain response to DTO
	dtoResp := dto.AuthResponse{
		APIToken: resp.APIToken,
		Username: req.Username,
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response{
		Success: true,
		Data:    dtoResp,
	})
}

// handleLogin handles user login
// @Summary Login user
// @Description Authenticate user with username and password to obtain an API token for protected endpoints
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.AuthResponse} "Successfully logged in"
// @Failure 400 {object} response.Response{error=response.APIError} "Invalid request or validation error"
// @Failure 401 {object} response.Response{error=response.APIError} "Invalid credentials"
// @Router /auth/login [post]
func (r *AuthRoutes) handleLogin(c *fiber.Ctx) error {
	var req dto.LoginRequest
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
	domainReq := domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := r.authService.Login(c.Context(), domainReq)
	if err != nil {
		logger.Error("login failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error: &response.APIError{
				Type:    "Unauthorized",
				Message: "Invalid credentials",
			},
		})
	}

	// Convert domain response to DTO
	dtoResp := dto.AuthResponse{
		APIToken: resp.APIToken,
		Username: req.Username,
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Success: true,
		Data:    dtoResp,
	})
}
