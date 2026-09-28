package server

import (
	"strconv"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/dto"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
// @Router /cart [get]
func (s *Server) getCart(c *gin.Context) {
	UserId := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(UserId)
	if err != nil {
		utils.NotFoundResponse(c, "cart not found")
		return
	}
	utils.SuccessResponse(c, "cart retrieved", cart)
}

// @Summary Add item to cart
// @Description Add a product to the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Item to add to cart"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items [post]
func (s *Server) addToCart(c *gin.Context) {
	UserID := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}
	cart, err := s.cartService.AddToCart(UserID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "failed to add item to cart", err)
		return
	}
	utils.SuccessResponse(c, "cart added to cart", cart)
}

// @Summary Update cart item quantity
// @Description Update the quantity of an item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{id} [put]
func (s *Server) updateCart(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "invalid id", err)
		return
	}
	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}
	cart, err := s.cartService.UpdateCartItem(UserID, uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, "failed to update item to cart", err)
		return
	}
	utils.SuccessResponse(c, "cart updated to cart", cart)
}

// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 400 {object} utils.Response "Invalid cart item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{id} [delete]
func (s *Server) removeFromCart(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "invalid id", err)
		return
	}

	if err := s.cartService.RemoveFromCart(UserID, uint(id)); err != nil {
		utils.BadRequestResponse(c, "failed to remove item from cart", err)
		return
	}
	utils.SuccessResponse(c, "cart removed from cart", nil)

}
