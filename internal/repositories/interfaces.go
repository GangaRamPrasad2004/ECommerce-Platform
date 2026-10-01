package repositories

import (
	
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	
)

// UserRepositoryInterface defines persistence operations for users and refresh tokens.
type UserRepositoryInterface interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetValidRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByID(id uint) error
}

// CartRepositoryInterface defines persistence operations for carts and cart items.
type CartRepositoryInterface interface {
	GetByUserID(userID uint) (*models.Cart, error)
	GetByUserIDWithItems(userID uint) (*models.Cart, error)
	GetOrCreateByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(id uint) error

	GetItemByCartAndProduct(cartID, productID uint) (*models.CartItem, error)
	GetItemByIDForUser(userID, itemID uint) (*models.CartItem, error)
	CreateItem(item *models.CartItem) error
	UpdateItem(item *models.CartItem) error
	DeleteItemForUser(userID, itemID uint) error
}

// ProductRepositoryInterface defines persistence operations for products and categories.
type ProductRepositoryInterface interface {
	// Category DB Operations
	CreateCategory(category *models.Category) error
	GetActiveCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error

	// Product DB Operations
	CreateProduct(product *models.Product) error
	GetProducts(offset, limit int) ([]models.Product, int64, error) // Returns items + total count
	GetProductByID(id uint) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error

	// Product Image DB Operations
	CountProductImages(productID uint) (int64, error)
	CreateProductImage(image *models.ProductImage) error

	// Search DB Operations
	SearchProducts(query string, categoryID *uint, minPrice, maxPrice *float64, offset, limit int) ([]ProductWithRank, int64, error)
}


// OrderRepositoryInterface defines persistence operations for customer orders.
type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	CreateFromCart(userID uint) (*models.Order, error)
	GetByUserID(userID uint, offset, limit int) ([]models.Order, int64, error)
	GetByIDForUser(userID, orderID uint) (*models.Order, error)
}
