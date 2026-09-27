package server

import (
	"strconv"

	"github.com/GangaRamPrasad2004/learning-go-shop/internal/dto"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(c *gin.Context) {
	UserId := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(UserId)
	if err != nil {
		utils.NotFoundResponse(c, "cart not found")
		return
	}
	utils.SuccessResponse(c, "cart retrieved", cart)
}

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

func (s *Server) removeFromCart(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "invalid id", err)
		return
	}

	if err := s.cartService.DeleteCartItem(UserID, uint(id)); err != nil {
		utils.BadRequestResponse(c, "failed to remove item from cart", err)
		return
	}
	utils.SuccessResponse(c, "cart removed from cart", nil)

}
