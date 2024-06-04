package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"payment-system-one/internal/api"
	"payment-system-one/internal/middleware"
	"payment-system-one/internal/ports"
	"time"
)

// SetupRouter is where router endpoints are called
func SetupRouter(handler *api.HTTPHandler, repository ports.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("/")
	{
		r.GET("/", handler.Readiness)
		r.POST("/create", handler.CreateAdmin)
		r.POST("/login", handler.LoginUer)
	}

	// authorizeAdmin authorizes all authorized Admins handlers
	authorizeAdmin := r.Group("/admin")
	authorizeAdmin.Use(middleware.AuthorizeAdmin(repository.FindAdminByEmail, repository.TokenInBlacklist))
	{
		authorizeAdmin.GET("/Admin", handler.GetAdminByEmail)
	}

	return router
}
