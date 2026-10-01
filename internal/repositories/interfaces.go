package repositories

import (
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
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
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error

	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error)
	GetProduct(id uint) (*dto.ProductResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error

	AddProductImage(productID uint, url, altText string) error
	SearchProducts(req *dto.SearchProductsRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error)
}

// OrderRepositoryInterface defines persistence operations for customer orders.
type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	CreateFromCart(userID uint) (*models.Order, error)
	GetByUserID(userID uint, offset, limit int) ([]models.Order, int64, error)
	GetByIDForUser(userID, orderID uint) (*models.Order, error)
}
