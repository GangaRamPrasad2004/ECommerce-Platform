package server

import (
	"net/http"

	_ "github.com/GangaRamPrasad2004/ECommerce-Platform/docs"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/config"
	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Server manages the application's HTTP server dependencies and routes.
type Server struct {
	config         *config.Config //nolint:gofmt // Field alignment is intentional for readability and existing config shape.
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *service.AuthService
	userService    *service.UserService
	productService *service.ProductService
	uploadService  *service.UploadService
	cartService    *service.CartService
	orderService   *service.OrderService
}

// New creates a Server with its shared configuration, database, logger, and services.
func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger, authService *service.AuthService, userService *service.UserService, productService *service.ProductService, uploadService *service.UploadService, cartService *service.CartService, orderService *service.OrderService) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		userService:    userService,
		uploadService:  uploadService,
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
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

	// documentation routes
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	router.GET("/playground", s.playgroundHandler())
	router.GET("/playground/public", s.playgroundPublicHandler())
	router.GET("/playground/protected", s.playgroundProtectedHandler())

	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(s.graphqlMiddleware())
	graphqlPublic.POST("/", s.graphqlHandler())

	graphqlProtected := router.Group("/graphql")
	graphqlProtected.Use(s.authMiddleware())
	graphqlProtected.Use(s.graphqlMiddleware())
	graphqlProtected.POST("/", s.graphqlHandler())

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)
	auth.POST("/refresh", s.refreshToken)
	auth.POST("/logout", s.logout)

	protected := api.Group("/")
	protected.Use(s.authMiddleware())

	// User routes
	users := protected.Group("/users")
	users.GET("/profile", s.getProfile)
	users.PUT("/profile", s.updateProfile)

	categories := protected.Group("/categories")
	categories.POST("/", s.adminMiddleware(), s.createCategory)
	categories.PUT("/:id", s.adminMiddleware(), s.updateCategory)
	categories.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)

	// product routes
	products := protected.Group("/products")
	products.POST("/", s.adminMiddleware(), s.createProduct)
	products.PUT("/:id", s.adminMiddleware(), s.updateProduct)
	products.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
	products.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)

	// cart routes
	cart := protected.Group("/carts")
	cart.GET("/", s.getCart)
	cart.POST("/items", s.addToCart)
	cart.PUT("/items/:id", s.updateCart)
	cart.DELETE("/items/:id", s.removeFromCart)

	// order routes
	orders := protected.Group("/orders")
	orders.POST("/", s.createOrder)
	orders.GET("/", s.getOrders)
	orders.GET("/:id", s.getOrder)

	// public routes
	api.GET("/categories", s.getCategories)
	api.GET("/search", s.searchProducts)
	api.GET("/products", s.getProducts)
	api.GET("/products/:id", s.getProduct)

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
