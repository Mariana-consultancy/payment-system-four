package api

import (
	"net/http"
	"os"
	"payment-system-one/internal/middleware"
	"payment-system-one/internal/models"
	"payment-system-one/internal/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) CreateAdmin(c *gin.Context) {
	var admin *models.Admin
	if err := c.ShouldBind(&admin); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//validate admin email

	//validate admin password

	//persist information in the data base
	err := u.Repository.CreateAdmin(admin)
	if err != nil {
		util.Response(c, "admin not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "admin created", 200, "success", nil)
}

// login
func (u *HTTPHandler) LoginAdmin(c *gin.Context) {
	var loginRequest *models.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}
	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Please enter your email or password", 400, "bad request body", nil)
		return
	}

	admin, err := u.Repository.FindAdminByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "admin does not exist", 404, "admin not found", nil)
		return
	}
	if admin.LoginCounter >= 3 {
		admin.IsLocked = true
		admin.UpdatedAt = time.Now()
		err = u.Repository.UpdateAdmin(admin)
		if err != nil {
			return
		}
		util.Response(c, "Your account has been lock after 3 failed attempt, contact customer care for assistance", 200, "success", nil)
		return
	}

	if admin.Password != loginRequest.Password {
		admin.LoginCounter++
		err := u.Repository.UpdateAdmin(admin)
		if err != nil {
			util.Response(c, "internal server error", 500, "admin not found", nil)
			return
		}
		util.Response(c, "password mismatch", 404, "admin not found", nil)
		return
	}

	//Generate token
	accessClaims, refreshClaims := middleware.GenerateClaims(admin.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "error generating access token", 500, "error generating access token", nil)
		return
	}
	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "error generating refresh token", 500, "error generating refresh token", nil)
		return
	}
	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "login successful", http.StatusOK, gin.H{
		"admin":         admin,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// call a protected route
func (u *HTTPHandler) GetAdminByEmail(c *gin.Context) {
	_, err := u.GetAdminFromContext(c)
	if err != nil {
		util.Response(c, "Admin not logged in", 500, "Admin not found", nil)
		return
	}

	email := c.Query("email")

	if email == "" {
		util.Response(c, "email is required", 400, "email is required", nil)
		return
	}

	admin, err := u.Repository.FindAdminByEmail(email)
	if err != nil {
		util.Response(c, "admin not fount", 500, "admin not found", nil)
		return
	}

	util.Response(c, "admin found", 200, admin, nil)
}
