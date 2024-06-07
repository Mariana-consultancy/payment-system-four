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

// Create a user
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	// if user.Password != user.ValidatePassword {
	// 	util.Response(c, "passwords do not match", 400, "bad request body", nil)
	// 	return
	// }

	//validate user email
	if !util.IsValidEmail(user.Email) {
		util.Response(c, "Invalid email address", 400, "Bad request body", nil)
		return
	}

	//check if user already exists
	_, err := u.Repository.FindUserByEmail(user.Email)
	if err == nil {
		util.Response(c, "User already exists", 400, "Bad request body", nil)
		return
	}

	//hash (hide) user password
	hashPass, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Could not hash password", 500, "internal server error", nil)
		return
	}
	user.Password = hashPass

	//generate account number
	acctNo, err := util.GenerateAccountNumber()
	if err != nil {
		util.Response(c, "Could not generate account number", 500, "internal server error", nil)
		return
	}

	user.AccountNumber = acctNo

	//set available balance to 0
	user.AvailableBalance = 500.0

	//persist information in the data base
	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "user not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "user created", 200, "success", nil)
}

// login
func (u *HTTPHandler) LoginUer(c *gin.Context) {
	var loginRequest *models.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}
	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Please enter your email or password", 400, "bad request body", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "user does not exist", 404, "user not found", nil)
		return
	}
	if user.LoginCounter >= 3 {
		user.IsLocked = true
		user.UpdatedAt = time.Now()
		err = u.Repository.UpdateUser(user)
		if err != nil {
			return
		}
		util.Response(c, "Your account has been locked after 3 failed attempts, contact customer care for assistance", 200, "success", nil)
		return
	}

	if user.Password != loginRequest.Password {
		user.LoginCounter++
		err := u.Repository.UpdateUser(user)
		if err != nil {
			util.Response(c, "internal server error", 500, "user not found", nil)
			return
		}
		util.Response(c, "password mismatch", 404, "user not found", nil)
		return
	}

	//Generate token
	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

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
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// call a protected route
func (u *HTTPHandler) GetUserByEmail(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 500, "user not found", nil)
		return
	}

	email := c.Query("email")

	if email == "" {
		util.Response(c, "email is required", 400, "email is required", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(email)
	if err != nil {
		util.Response(c, "user not found", 500, "user not found", nil)
		return
	}

	util.Response(c, "user found", 200, user, nil)
}

//query parameter
//path parameter

//100 ---- informtional
//200 ---- success 200, 201, 202
//300 ---- redirect
//400 ---- client error
//500 ----- server error

//syntax error
