package service

import (
	"mime/multipart"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
)

// AuthServiceInterface defines the contract for authentication-related operations.
type AuthServiceInterface interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
}

// UserServiceInterface defines the contract for user profile operations.
type UserServiceInterface interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

// ProductServiceInterface defines the contract for product and category management operations.
type ProductServiceInterface interface {
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

// CartServiceInterface defines the contract for cart-related operations.
type CartServiceInterface interface {
	GetCart(userID uint) (*dto.CartResponse, error)
	AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveFromCart(userID, itemID uint) error
}

// OrderServiceInterface defines the contract for order-related operations.
type OrderServiceInterface interface {
	CreateOrder(userID uint) (*dto.OrderResponse, error)
	GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error)
	GetOrder(userID, orderID uint) (*dto.OrderResponse, error)
}

// UploadServiceInterface defines the contract for product image upload operations.
type UploadServiceInterface interface {
	UploadProductImage(productID uint, file *multipart.FileHeader) (string, error)
}
