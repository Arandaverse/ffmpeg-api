package handlers

import (
	"ffmpeg-api/internal/handlers/routes"
	"ffmpeg-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// Handler holds all HTTP handlers and their dependencies
type Handler struct {
	authRoutes   *routes.AuthRoutes
	ffmpegRoutes *routes.FFMPEGRoutes
	indexRoutes  *routes.IndexRoutes
}

// NewHandler creates a new Handler instance
func NewHandler(authService service.AuthService, ffmpegService service.FFMPEGService) *Handler {
	return &Handler{
		authRoutes:   routes.NewAuthRoutes(authService),
		ffmpegRoutes: routes.NewFFMPEGRoutes(ffmpegService, authService),
		indexRoutes:  routes.NewIndexRoutes(),
	}
}

// RegisterRoutes registers all routes to the given router
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Register index routes first
	h.indexRoutes.Register(app)

	// Register auth routes
	h.authRoutes.Register(app)

	// Register FFMPEG routes
	h.ffmpegRoutes.Register(app)
}

// ErrorHandler handles errors returned from routes
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"type":    fiber.ErrInternalServerError.Error(),
			"message": err.Error(),
		},
	})
}

// NewFiberApp creates a new Fiber app with configured template engine
func NewFiberApp() *fiber.App {
	// Create a new template engine
	engine := html.New("./views", ".html")

	// Create a new Fiber app with the template engine
	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: ErrorHandler,
	})

	return app
}
