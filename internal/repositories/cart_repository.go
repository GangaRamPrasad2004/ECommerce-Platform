package repositories

import (
	"errors"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

var _ CartRepositoryInterface = (*CartRepository)(nil)

// CartRepository stores and retrieves carts and their items.
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a cart repository backed by db.
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

// Create stores a cart.
func (repo *CartRepository) Create(cart *models.Cart) error {
	return repo.db.Create(cart).Error
}

// Update saves changes to a cart.
func (repo *CartRepository) Update(cart *models.Cart) error {
	return repo.db.Save(cart).Error

}

// Delete removes a cart by ID.
func (repo *CartRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Cart{}, id).Error
}

// GetByUserID retrieves a user's cart without preloaded items.
func (repo *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := repo.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

// GetByUserIDWithItems retrieves a user's cart with its items, products, and categories.
func (repo *CartRepository) GetByUserIDWithItems(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := repo.db.Preload("CartItems.Product.Category").
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetOrCreateByUserID retrieves a user's cart or creates one if it does not exist.
func (repo *CartRepository) GetOrCreateByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := repo.db.Where("user_id = ?", userID).
		FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetItemByCartAndProduct retrieves a cart item, returning nil when it does not exist.
func (repo *CartRepository) GetItemByCartAndProduct(cartID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := repo.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetItemByIDForUser retrieves a cart item only if it belongs to the specified user.
func (repo *CartRepository) GetItemByIDForUser(userID, itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	if err := repo.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem stores a cart item.
func (repo *CartRepository) CreateItem(item *models.CartItem) error {
	return repo.db.Create(item).Error
}

// UpdateItem saves changes to a cart item.
func (repo *CartRepository) UpdateItem(item *models.CartItem) error {
	return repo.db.Save(item).Error
}

// DeleteItemForUser removes a cart item if it belongs to the specified user.
func (repo *CartRepository) DeleteItemForUser(userID, itemID uint) error {
	return repo.db.Where("id = ? AND cart_id IN (?)", itemID,
		repo.db.Select("id").Table("carts").Where("user_id = ?", userID)).
		Delete(&models.CartItem{}).Error
}
