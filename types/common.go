package types

type ComponentOptions struct {
	Display bool
	Profile string
}

// Define the context key type.
type contextKey string

// Create a context key for the theme.
var ProfileContextKey contextKey = "profile"
var Profile string = "guest"
