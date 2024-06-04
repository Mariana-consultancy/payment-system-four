package models

import "gorm.io/gorm"

/*type Admin struct {
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
}*/

type Admin struct {
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


//type AdminProfile struct {
//	gorm.Model
//	ValidIdentity string `json:"valid_identity"`
//	PassPort string `json:"passport"`
//
//}

type Transaction struct {
	gorm.Model
	AdminID          uint    `json:"admin_id"`
	Amount          float64 `json:"amount"`
	Reference       string  `json:"reference"`
	TransactionType string  `json:"transaction_type"`
	AccountNumber   string  `json:"account_number"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
