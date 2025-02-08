package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/repository"
	"fmt"
)

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Generate API token
	apiToken, err := s.generateAPIToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API token: %w", err)
	}

	// Create user
	user := &domain.User{
		Username: req.Username,
		APIToken: apiToken,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.AuthResponse{
		APIToken: apiToken,
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// In a real application, you would verify the password here
	// For this example, we just return the API token

	return &domain.AuthResponse{
		APIToken: user.APIToken,
	}, nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.userRepo.FindByAPIToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return user, nil
}

func (s *AuthServiceImpl) generateAPIToken() (string, error) {
	b := make([]byte, s.config.Server.APITokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
