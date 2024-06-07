package api

import (
	"payment-system-one/internal/util"

	"payment-system-one/internal/models"

	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) TransferFunds(c *gin.Context) {
	var funds *models.TransferFunds
	if err := c.ShouldBind(&funds); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	// declare request body

	//bind JSON data to struct

	//Get user from context - ensures the person who is logged in is the correct authorised person trying to transfer funds
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged", 500, "User not found", nil)
		return
	}

	//validate the amount you are trying to send (ensure it is not a number below 0)
	if funds.Amount <= 0 {
		util.Response(c, "Invalid amount", 500, "Bad request body", nil)
		return
	}

	//check if the account number exists
	recipient, err := u.Repository.FindUserByAccountNumber(funds.AccountNumber)
	if err != nil {
		util.Response(c, "User not found", 400, "Bad request body", nil)
		return
	}

	//make sure the trasnfer amount is less than or equal to the user's current balance
	if funds.Amount >= user.AvailableBalance {
		util.Response(c, "Insufficient funds", 400, "Bad request", nil)
		return
	}

	//persist the data into the DB
	err = u.Repository.TransferFunds(user, recipient, funds.Amount)
	if err != nil {
		util.Response(c, "Transfer not possible", 500, "Transfer not successful", nil)
		return
	}
	util.Response(c, "Transfer was done successfully", 200, "Transfer was successful", nil)

}
