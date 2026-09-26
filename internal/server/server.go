package server

import (
	"net/http"

	"github.com/GangaRamPrasad2004/learning-go-shop/internal/config"
	"github.com/GangaRamPrasad2004/learning-go-shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Server manages the application's HTTP server dependencies and routes.
type Server struct {
	config         *config.Config //nolint:gofmt
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *service.AuthService
	userService    *service.UserService
	productService *service.ProductService
	uploadService  *service.UploadService
}

// New creates a Server with its shared configuration, database, logger, and services.
func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger, authService *service.AuthService, userService *service.UserService, productService *service.ProductService, uploadService *service.UploadService) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		userService:    userService,
		uploadService:  uploadService,
		productService: productService,
	}
}

// SetupRoutes configures and returns the application's HTTP router.
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()
	// add MiddleWares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// add routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running",
		})
	})
	router.GET("/health", s.healthCheck)
	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{ // auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}
			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			// product routes
			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
			}
		}
		//public routes
		api.GET("/categories", s.getCategories)
		//api.GET("/search", s.searchProducts)
		//api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)
	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}

}
