package auth

// User repo interface
type Repository interface {
	Create() error
}
