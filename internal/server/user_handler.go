package server

import (
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/dto"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/service"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getProfile(c *gin.Context) {

	userID := c.GetUint("user_id")
	userService := service.NewUserService(s.db)
	profile, err := userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}
func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	userService := service.NewUserService(s.db)
	profile, err := userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", profile)
}
