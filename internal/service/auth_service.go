package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"ffmpeg-api/internal/config"
	"ffmpeg-api/internal/domain"
	"ffmpeg-api/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

// Register registers a new user
func (s *AuthServiceImpl) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Generate API token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate API token: %w", err)
	}
	apiToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		APIToken: apiToken,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.AuthResponse{
		APIToken: apiToken,
	}, nil
}

// Login authenticates a user and returns an API token
func (s *AuthServiceImpl) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByUsernameWithPassword(ctx, req.Username)

	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &domain.AuthResponse{
		APIToken: user.APIToken,
	}, nil
}

// ValidateToken validates an API token and returns the associated user
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.userRepo.FindByAPIToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid API token")
	}
	return user, nil
}
