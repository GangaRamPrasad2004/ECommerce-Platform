package repositories

import (
	"time"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

// UserRepository stores and retrieves users and refresh tokens.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a user repository backed by db.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByEmail retrieves a user by email address.
func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID.
func (repo *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmailAndActive retrieves a user by email and active status.
func (repo *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ? AND is_active = ?", email, isActive).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create stores a user.
func (repo *UserRepository) Create(user *models.User) error {
	return repo.db.Create(user).Error
}

// Update saves changes to a user.
func (repo *UserRepository) Update(user *models.User) error {
	return repo.db.Save(user).Error
}

// Delete removes a user by ID.
func (repo *UserRepository) Delete(id uint) error {
	return repo.db.Delete(&models.User{}, id).Error
}

// CreateRefreshToken stores a refresh token.
func (repo *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return repo.db.Create(token).Error
}

// GetValidRefreshToken retrieves an unexpired refresh token.
func (repo *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := repo.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken removes a refresh token by its token value.
func (repo *UserRepository) DeleteRefreshToken(token string) error {
	return repo.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteRefreshTokenByID removes a refresh token by ID.
func (repo *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return repo.db.Delete(&models.RefreshToken{}, id).Error
}
