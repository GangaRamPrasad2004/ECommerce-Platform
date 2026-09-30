package resolver

import (
	"strconv"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds the services used by the GraphQL resolvers.
type Resolver struct {
	authService    service.AuthServiceInterface
	userService    service.UserServiceInterface
	productService service.ProductServiceInterface
	cartService    service.CartServiceInterface
	orderService   service.OrderServiceInterface
}

// NewResolver creates a new GraphQL resolver with all required services.
func NewResolver(authService service.AuthServiceInterface,
	userService service.UserServiceInterface,
	productService service.ProductServiceInterface,
	cartService service.CartServiceInterface,
	orderService service.OrderServiceInterface) *Resolver {

	return &Resolver{
		authService:    authService,
		userService:    userService,
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
	}

}

func (r *Resolver) parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	return uint(parsed), err
}
