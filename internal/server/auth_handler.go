package server

import (
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/dto"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/service"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid Request Data", err)
		return
	}
	authService := service.NewAuthService(s.db, s.config)
	response, err := authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "registration failed ", err)
		return
	}
	utils.CreatedResponse(c, "user Registratered Successfully", response)

}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid Request Data", err)
		return
	}
	authService := service.NewAuthService(s.db, s.config)
	response, err := authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Loin failed")
		return
	}
	utils.SuccessResponse(c, "Login Successful", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := service.NewAuthService(s.db, s.config)
	response, err := authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := service.NewAuthService(s.db, s.config)
	if err := authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, "Logout successful", nil)
}
