package service

import (
	"errors"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/repositories"
)

var _ CartServiceInterface = (*CartService)(nil)

// CartService manages shopping carts and their items.
type CartService struct {
	cartRepo    repositories.CartRepositoryInterface
	productRepo repositories.ProductRepositoryInterface
}

// NewCartService creates a cart service backed by cart and product repositories.
func NewCartService(cartRepo repositories.CartRepositoryInterface, productRepo repositories.ProductRepositoryInterface) *CartService {
	return &CartService{cartRepo: cartRepo, productRepo: productRepo}
}

// GetCart returns the cart belonging to the specified user.
func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetByUserIDWithItems(userID)
	if err != nil {
		return nil, err
	}
	return s.convertToCartResponse(cart), nil
}

// AddToCart adds a product quantity to the user's cart.
func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	product, err := s.productRepo.GetProductByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("stock is less or insufficient")
	}

	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, err
	}

	cartItem, err := s.cartRepo.GetItemByCartAndProduct(cart.ID, product.ID)
	if err != nil {
		return nil, err
	}
	if cartItem == nil {
		item := models.CartItem{
			ProductID: req.ProductID,
			CartID:    cart.ID,
			Quantity:  req.Quantity,
		}
		if err := s.cartRepo.CreateItem(&item); err != nil {
			return nil, err
		}
	} else {
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("insufficient quantity")
		}
		if err := s.cartRepo.UpdateItem(cartItem); err != nil {
			return nil, err
		}
	}
	return s.GetCart(userID)
}

// UpdateCartItem changes the quantity of an item in the user's cart.
func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cartItem, err := s.cartRepo.GetItemByIDForUser(userID, itemID)
	if err != nil {
		return nil, errors.New("CartItem not found")
	}
	product, err := s.productRepo.GetProductByID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("stock is less or insufficient")
	}
	cartItem.Quantity = req.Quantity
	if err := s.cartRepo.UpdateItem(cartItem); err != nil {
		return nil, err
	}
	return s.GetCart(userID)
}

// RemoveFromCart removes an item from the user's cart.
func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	return s.cartRepo.DeleteItemForUser(userID, itemID)
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64
	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}
	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
