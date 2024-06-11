package ports

import "payment-system-four/internal/models"

// Repository defines the interface for interacting with data storage.
type Repository interface {
	// FindUserByEmail retrieves a user by their email address.
	FindUserByEmail(email string) (*models.User, error)

	// TokenInBlacklist checks if a token is in the blacklist.
	TokenInBlacklist(token *string) bool

	// CreateUser creates a new user in the database.
	CreateUser(user *models.User) error

	// UpdateUser updates an existing user in the database.
	UpdateUser(user *models.User) error
	FindUserByAccountNumber(account_number int) (*models.User, error)
	TransferFunds(user *models.User, recipient *models.User, amount float64) error
}
