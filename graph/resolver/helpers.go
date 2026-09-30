package resolver

import (
	"context"
	"errors"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/utils"
)

var (
	// ErrUnauthorized indicates a missing or invalid user context in the GraphQL request.
	ErrUnauthorized = errors.New("unauthorized")
)

const (
	adminRole = "admin"
)

// GetUserIDFromContext extracts the authenticated user ID from the GraphQL context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value(utils.UserIDKey)
	if userID == nil {
		return 0, ErrUnauthorized
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, ErrUnauthorized
}

// GetUserRoleFromContext extracts the authenticated user role from the GraphQL context.
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value(utils.UserRoleKey)
	if userRole == nil {
		return "", ErrUnauthorized
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", ErrUnauthorized
}

// IsAdminFromContext reports whether the current request has the admin role.
func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}

func getPagingNumbers(page, limit *int) (pageNumber, pageLimit int) {
	var p, l = 0, 0

	if page != nil {
		p = *page
	}

	if limit != nil {
		l = *limit
	}

	if p <= 0 {
		p = 1
	}

	if l <= 0 {
		l = 10
	}

	return p, l
}
