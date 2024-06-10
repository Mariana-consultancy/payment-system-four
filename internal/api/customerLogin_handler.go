package api

import (
	
	"net/http"
	"os"
	

	//"os/user"
	"payment-system-one/internal/middleware"
	"payment-system-one/internal/models"
	"payment-system-one/internal/util"


	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//create a user
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//validate admin email
	if !util.IsValidEmail(user.Email){
		util.Response(c,"invail email",400,"Bad requst body",nil)
		return
	}
	//check if the user already exist
	_, err := u.Repository.FindUserByEmail(user.Email)
	if err == nil {
		util.Response(c, "User already exist",400,"Bad requst body",nil)
		return
	}
	hashPass, err := util.HashPassword(user.Password)
	if err !=nil {
		util.Response(c,"could not harsh password",500,"internal server error",nil)
		return
	}
	user.Password = hashPass

	//validate admin password

	//persist information in the data base
	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "user not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "user created", 200, "success", nil)
}


// login
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}
	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Please enter your email and/or password", 400, "bad request body", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "user does not exist", 404, "user not found", nil)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password));err != nil{
		util.Response(c, "invalid email or passowrd", 400, "invalid email or paswword",nil)
		return
	}
	//if user.LoginCounter >= 3 {
		//user.IsLocked = true
		//user.UpdatedAt = time.Now()
		//err = u.Repository.UpdateUser(user)
		//if err != nil {
		//	return
		//}
		//util.Response(c, "Your account has been lock after 3 failed attempt, contact customer care for assistance", 200, "success", nil)
		//return
	//}

	//if user.Password != loginRequest.Password {
	//	user.LoginCounter++
	//	err := u.Repository.UpdateUser(user)
	//	if err != nil {
	//		util.Response(c, "internal server error", 500, "user not found", nil)
	//		return
	//	}
	//	util.Response(c, "password mismatch", 404, "user not found", nil)
	//	return
	//}



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
		"user":         user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// call a protected route
func (u *HTTPHandler) GetUserByEmail(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "user not logged in", 401, "user not found", nil)
		return
	}

	email := c.Query("email")

	if email == "" {
		util.Response(c, "email is required", 400, "email is required", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(email)
	if err != nil {
		util.Response(c, "user not fount", 400, "user not found", nil)
		return
	}

	util.Response(c, "user found", 200, user, nil)
}
