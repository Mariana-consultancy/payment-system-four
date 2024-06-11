package models

import "gorm.io/gorm"

// User struct represents a user in the system.
type User struct {
	gorm.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Password     string `json:"password"`
	DateOfBirth  string `json:"date_of_birth"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	LoginCounter int    `json:"login_counter"`
	IsLocked     bool   `json:"is_locked"`
}

// Uncomment and expand the UserProfile struct if additional user profile details are needed.
// type UserProfile struct {
// 	gorm.Model
// 	ValidIdentity string `json:"valid_identity"` // Valid identity proof for the user
// 	PassPort      string `json:"passport"`       // User's passport information
// }

// Transaction struct represents a financial transaction made by a user.
type Transaction struct {
	gorm.Model
	UserID          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	Reference       string  `json:"reference"`
	TransactionType string  `json:"transaction_type"`
}

// LoginRequest struct is used for handling login requests.
type LoginRequest struct {
	Email    string `json:"email"`    // Email address provided for login
	Password string `json:"password"` // Password provided for login
}
