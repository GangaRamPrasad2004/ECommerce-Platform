package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/config"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/events"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
	"gorm.io/gorm"
)

var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService handles user registration, authentication, and token lifecycle operations.
type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.Publisher
}

// NewAuthService creates an AuthService backed by db and configured with config.
func NewAuthService(db *gorm.DB, cfg *config.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         cfg,
		eventPublisher: eventPublisher,
	}

}

// Register creates a customer account and returns its authentication tokens.
func (s *AuthService) Register(regReq *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser *models.User
	err := s.db.Where("email = ?", regReq.Email).First(&existingUser).Error

	if err == nil {
		return nil, errors.New("user already exists with this email")
	}

	hashedPassword, err := utils.HashPassword(regReq.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Email:     regReq.Email,
		Password:  hashedPassword,
		FirstName: regReq.FirstName,
		LastName:  regReq.LastName,
		Phone:     regReq.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("unable to create the cart")
	}
	return s.generatedAuthResponse(&user)
}

// Login validates a user's credentials and returns its authentication tokens.
func (s *AuthService) Login(logReq *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", logReq.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid Credentials") // Security Best Practice
	}
	if !utils.CheckPassword(logReq.Password, user.Password) {
		return nil, errors.New("invalid Credentials") // Security Best Practice
	}
	return s.generatedAuthResponse(&user)
}

// RefreshToken validates and replaces a refresh token with a new token pair.
func (s *AuthService) RefreshToken(refReq *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(refReq.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	var refreshToken models.RefreshToken
	if err := s.db.Where("token= ? AND expires_at > ?", refReq.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New(" refresh Token not found ")
	}
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	s.db.Delete(&refreshToken)
	return s.generatedAuthResponse(&user)
}

// Logout invalidates the supplied refresh token.
func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}
func (s *AuthService) generatedAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}
	s.db.Create(&refreshTokenModel)

	err = s.eventPublisher.Publish("UserLoggedIn", user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish user login event: %w", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
