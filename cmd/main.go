package main

import (
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/server"
)

// @title FFMPEG Serverless API
// @version 1.0
// @description A serverless API for processing videos using FFMPEG. This API allows you to submit video processing jobs, monitor their progress, and manage user authentication.

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Token
// @description API token obtained after login. Required for all protected endpoints.

// @tag.name Auth
// @tag.description Authentication endpoints for user registration and login

// @tag.name FFMPEG
// @tag.description Video processing endpoints using FFMPEG

// @tag.name Index
// @tag.description Main page and general information

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
