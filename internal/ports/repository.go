package ports

import "payment-system-one/internal/models"

type Repository interface {
	FindAdminByEmail(email string) (*models.Admin, error)
	TokenInBlacklist(token *string) bool
	CreateAdmin(Admin *models.Admin) error
	UpdateAdmin(Admin *models.Admin) error
}
