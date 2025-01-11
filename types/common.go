package types

type ComponentOptions struct {
	Display bool
	Role    string
}

// API response
type Response struct {
	Status  string
	Message string
}

type User struct {
	Uuid           string
	Name           string
	Email          string
	Role           string
	CreatedAt      string
	UpdatedAt      string
	GeneratedPhoto string
}
