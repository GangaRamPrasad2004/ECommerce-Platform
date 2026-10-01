package repositories

import (
	"errors"
	"fmt"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/models"
	"gorm.io/gorm"
)

var _ OrderRepositoryInterface = (*OrderRepository)(nil)

// OrderRepository stores and retrieves customer orders.
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates an order repository backed by db.
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create stores an order.
func (repo *OrderRepository) Create(order *models.Order) error {
	return repo.db.Create(order).Error
}

// CreateFromCart creates an order from a user's cart and updates stock atomically.
func (repo *OrderRepository) CreateFromCart(userID uint) (*models.Order, error) {
	var createdOrder models.Order
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return errors.New("cart not found")
		}
		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}

		var totalAmount float64
		orderItems := make([]models.OrderItem, 0, len(cart.CartItems))
		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]
			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", cartItem.Product.Name)
			}

			totalAmount += float64(cartItem.Quantity) * cartItem.Product.Price
			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})

			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(&cartItem.Product).Error; err != nil {
				return err
			}
		}

		order := models.Order{
			UserID:      userID,
			Status:      models.OrderStatusPending,
			TotalAmount: totalAmount,
			OrderItems:  orderItems,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		if err := tx.Preload("OrderItems.Product.Category").First(&createdOrder, order.ID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &createdOrder, nil
}

// GetByUserID retrieves a page of a user's orders and the total order count.
func (repo *OrderRepository) GetByUserID(userID uint, offset, limit int) ([]models.Order, int64, error) {
	var total int64
	if err := repo.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	if err := repo.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetByIDForUser retrieves an order only if it belongs to the specified user.
func (repo *OrderRepository) GetByIDForUser(userID, orderID uint) (*models.Order, error) {
	var order models.Order
	if err := repo.db.Preload("OrderItems.Product.Category").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
