package server

import (
	_ "ffmpeg-api/docs"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/database"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/handlers"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/repository"
	"ffmpeg-api/internal/service"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	db     database.Database
}

func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&domain.User{}, &domain.JobStatus{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Create repositories
	userRepo := repository.NewGormUserRepository(db)
	jobRepo := repository.NewGormJobRepository(db)

	// Create storage service based on configuration
	storageService, err := initStorageService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage service: %w", err)
	}

	// Create services
	authService := service.NewAuthService(userRepo, cfg)
	ffmpegService := service.NewFFMPEGService(jobRepo, userRepo, storageService, cfg)

	// Create Fiber app
	app := handlers.NewFiberApp()

	// Add middlewares
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberLogger.New())

	// Create handlers
	handler := handlers.NewHandler(authService, ffmpegService)

	// Swagger documentation
	app.Get("/swagger/*", swagger.New(swagger.Config{
		PersistAuthorization: true,
	}))

	// Register routes
	handler.RegisterRoutes(app)

	return &Server{
		app:    app,
		config: cfg,
		db:     db,
	}, nil
}

func (s *Server) Start() error {
	// Create temp directories
	createTempDirectories(s.config)

	// Start server
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	logger.Info("server starting", "address", addr)
	return s.app.Listen(addr)
}

func (s *Server) Close() error {
	return s.db.Close()
}

func initStorageService(cfg *config.Config) (service.StorageService, error) {
	if cfg.Storage.Provider == "minio" {
		logger.Info("using MinIO storage service")
		return service.NewMinioStorageService(cfg)
	}
	logger.Info("using local storage service")
	return service.NewLocalStorageService(cfg), nil
}

func createTempDirectories(cfg *config.Config) {
	dirs := []string{
		cfg.FFMPEG.TempDirectory,
		cfg.Storage.TempDirectory,
		fmt.Sprintf("%s/uploads", cfg.Storage.TempDirectory),
	}

	for _, dir := range dirs {
		logger.Debug("creating directory", "path", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Error("failed to create directory", "path", dir, "error", err)
		}
	}
}
