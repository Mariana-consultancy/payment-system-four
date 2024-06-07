package server

import (
	"payment-system-one/internal/api"
	"payment-system-one/internal/middleware"
	"payment-system-one/internal/ports"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		r.POST("/create", handler.CreateUser)
		r.POST("/login", handler.LoginUer)
		r.POST("/funds/transfer", handler.TransferFunds)
	}

	// authorizeUser authorizes all authorized users handlers
	authorizeUser := r.Group("/user")
	authorizeUser.Use(middleware.AuthorizeAdmin(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		authorizeUser.POST("/funds/transfer", handler.TransferFunds)
	}

	// authorizeAdmin authorizes all authorized users handlers
	authorizeAdmin := r.Group("/admin")
	authorizeAdmin.Use(middleware.AuthorizeAdmin(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		authorizeAdmin.GET("/user", handler.GetUserByEmail)
	}

	return router
}
