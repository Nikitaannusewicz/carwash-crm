package middleware

type ContextKey string

const (
	RoleKey   ContextKey = "role"
	UserIDKey ContextKey = "userID"
)
