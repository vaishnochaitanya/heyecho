package types

type ComponentOptions struct {
	Display bool
	Role    string
}

// Define the context key type.
type contextKey string

// Create a context key for the theme.
var RoleContextKey contextKey = "role"
var Role string = "guest"
