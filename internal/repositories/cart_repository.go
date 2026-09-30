package repositories

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (repo *CartRepository) Create(cart *models.Cart) error {
	return repo.db.Create(cart).Error
}
func (repo *CartRepository) Update(cart *models.Cart) error {
	return repo.db.Save(cart).Error

}
func (repo *CartRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Cart{}, id).Error
}
func (repo *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := repo.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil,err
	}

	return &cart,nil
}
