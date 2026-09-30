package utils

// ContextKey identifies a value stored in the request context.
type ContextKey string

const (
	// UserIDKey stores the authenticated user ID in context.
	UserIDKey ContextKey = "user_id"
	// UserEmailKey stores the authenticated user email in context.
	UserEmailKey ContextKey = "user_email"
	// UserRoleKey stores the authenticated user role in context.
	UserRoleKey ContextKey = "user_role"
	// GinContextKey stores the Gin context in the request context.
	GinContextKey ContextKey = "gin_context"
)
