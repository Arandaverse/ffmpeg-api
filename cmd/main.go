package main

import (
	_ "ffmpeg-api/docs" // Import swagger docs
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/handlers"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/repository"
	"ffmpeg-api/internal/service"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title FFMPEG Serverless API
// @version 1.0
// @description A serverless API for processing videos using FFMPEG

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Token
// @description API token for authentication

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load configuration", "error", err)
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to initialize database", "error", err)
	}

	// Create repositories
	userRepo := repository.NewGormUserRepository(db)
	jobRepo := repository.NewGormJobRepository(db)

	// Create storage service based on configuration
	var storageService service.StorageService
	if cfg.Storage.Provider == "minio" {
		logger.Info("using MinIO storage service")
		storageService, err = service.NewMinioStorageService(cfg)
		if err != nil {
			logger.Fatal("failed to initialize MinIO storage service", "error", err)
		}
	} else {
		logger.Info("using local storage service")
		storageService = service.NewLocalStorageService(cfg)
	}

	// Create services
	authService := service.NewAuthService(userRepo, cfg)
	ffmpegService := service.NewFFMPEGService(jobRepo, userRepo, storageService, cfg)

	// Create Fiber app
	app := handlers.NewFiberApp()

	// Add middlewares
	// app.Use(handlers.LoggingMiddleware())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberLogger.New())

	// Create handlers
	handler := handlers.NewHandler(authService, ffmpegService)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes
	handler.RegisterRoutes(app)

	// Create temp directories if they don't exist
	createTempDirectories(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("server starting", "address", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal("server failed", "error", err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger.Info("initializing database", "driver", cfg.Database.Driver)
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	logger.Info("running database migrations")
	if err := db.AutoMigrate(&domain.User{}, &domain.JobStatus{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
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
