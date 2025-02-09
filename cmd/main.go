package main

import (
	_ "ffmpeg-api/docs" // Import swagger docs
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/server"
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

	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		logger.Fatal("failed to create server", "error", err)
	}

	if err := srv.Start(); err != nil {
		logger.Fatal("server failed", "error", err)
	}
}
