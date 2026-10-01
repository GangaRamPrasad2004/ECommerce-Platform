package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/config"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/events"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/repositories"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
)

var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService handles user registration, authentication, and token lifecycle operations.
type AuthService struct {
	userRepo       repositories.UserRepositoryInterface
	cartRepo       repositories.CartRepositoryInterface
	config         *config.Config
	eventPublisher events.Publisher
}

// NewAuthService creates an AuthService backed by db and configured with config.
func NewAuthService(cfg *config.Config,
	eventPublisher events.Publisher,
	userRepo repositories.UserRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface,
) *AuthService {
	return &AuthService{

		userRepo:       userRepo,
		cartRepo:       cartRepo,
		config:         cfg,
		eventPublisher: eventPublisher,
	}

}

// Register creates a customer account and returns its authentication tokens.
func (s *AuthService) Register(regReq *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.userRepo.GetByEmail(regReq.Email); err == nil {
		return nil, errors.New("you cannot register with this email")
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
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		fmt.Println("unable to create the cart")
	}
	return s.generatedAuthResponse(&user)
}

// Login validates a user's credentials and returns its authentication tokens.
func (s *AuthService) Login(logReq *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(logReq.Email, true)
	if err != nil {
		return nil, errors.New("invalid Credentials")
	}
	if !utils.CheckPassword(logReq.Password, user.Password) {
		return nil, errors.New("invalid Credentials")
	}
	return s.generatedAuthResponse(user)
}

// RefreshToken validates and replaces a refresh token with a new token pair.
func (s *AuthService) RefreshToken(refReq *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(refReq.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := s.userRepo.GetValidRefreshToken(refReq.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh Token not found")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.DeleteRefreshToken(refreshToken.Token); err != nil {
		return nil, err
	}
	return s.generatedAuthResponse(user)
}

// Logout invalidates the supplied refresh token.
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
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
	if err := s.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		return nil, err
	}

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
