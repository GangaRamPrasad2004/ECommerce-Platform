package repositories

import (
	"time"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
func (repo *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ? AND is_active = ?", email, isActive).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepository) Create(user *models.User) error {
	return repo.db.Create(user).Error
}
func (repo *UserRepository) Update(user *models.User) error {
	return repo.db.Save(user).Error
}
func (repo *UserRepository) Delete(id uint) error {
	return repo.db.Delete(&models.User{}, id).Error
}

func (repo *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return repo.db.Create(token).Error
}
func (repo *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := repo.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}
func (repo *UserRepository) DeleteRefreshToken(token string) error {
	return repo.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}
func (repo *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return repo.db.Delete(&models.RefreshToken{}, id).Error
}
