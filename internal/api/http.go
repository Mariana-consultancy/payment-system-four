package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"payment-system-one/internal/models"
	"payment-system-one/internal/ports"
)

type HTTPHandler struct {
	Repository ports.Repository
}

func NewHTTPHandler(repository ports.Repository) *HTTPHandler {
	return &HTTPHandler{
		Repository: repository,
	}
}

func (u *HTTPHandler) GetAdminFromContext(c *gin.Context) (*models.Admin, error) {
	contextAdmin, exists := c.Get("admin")
	if !exists {
		return nil, fmt.Errorf("error getting user from context")
	}
	admin, ok := contextAdmin.(*models.Admin)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return admin, nil
}

func (u *HTTPHandler) GetTokenFromContext(c *gin.Context) (string, error) {
	tokenI, exists := c.Get("access_token")
	if !exists {
		return "", fmt.Errorf("error getting access token")
	}
	tokenstr := tokenI.(string)
	return tokenstr, nil
}
