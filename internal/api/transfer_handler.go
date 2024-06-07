package api

import (
	"github.com/gin-gonic/gin"
	"payment-system-four/internal/models"
	"payment-system-four/internal/util"
)

// TransferHandler handles the transfer of funds from one user to another
func (u *HTTPHandler) TransferHandler(c *gin.Context) {
	var transferRequest struct {
		ReceiverName  string  `json:"receiver_name"`
		ReceiverEmail string  `json:"receiver_email"`
		Amount        float64 `json:"amount"`
	}

	// Bind the request body to the transferRequest struct
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		util.Response(c, "Invalid request", 400, "Bad request body", nil)
		return
	}

	// Get the sender user from the context
	senderUser, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Sender user not found", 404, "Sender user not found", nil)
		return
	}

	// Find the receiver user by email
	receiverUser, err := u.Repository.FindUserByEmail(transferRequest.ReceiverEmail)
	if err != nil {
		util.Response(c, "Receiver user not found", 404, "Receiver user not found", nil)
		return
	}

	// Check if the sender has sufficient balance
	if senderUser.InitialBalance < transferRequest.Amount {
		util.Response(c, "Insufficient balance", 400, "Insufficient balance", nil)
		return
	}

	// Create a new transaction for the transfer
	transaction := &models.Transaction{
		UserID:          senderUser.ID,
		Amount:          -transferRequest.Amount, // Negative amount for the sender
		TransactionType: "transfer",
		Reference:       receiverUser.Email, // Use the receiver's email as the reference
	}

	// Save the transaction to the database
	if err := u.Repository.CreateTransaction(transaction); err != nil {
		util.Response(c, "Error creating transaction", 500, "Internal server error", nil)
		return
	}

	// Create a new transaction for the receiver
	receiverTransaction := &models.Transaction{
		UserID:          receiverUser.ID,
		Amount:          transferRequest.Amount,
		TransactionType: "transfer",
		Reference:       senderUser.Email, // Use the sender's email as the reference
	}

	// Save the receiver transaction to the database
	if err := u.Repository.CreateTransaction(receiverTransaction); err != nil {
		util.Response(c, "Error creating transaction", 500, "Internal server error", nil)
		return
	}

	// Update the sender's balance
	senderUser.InitialBalance -= transferRequest.Amount

	// Update the receiver's balance
	receiverUser.InitialBalance += transferRequest.Amount

	// Save the updated users to the database
	if err := u.Repository.UpdateUser(senderUser); err != nil {
		util.Response(c, "Error updating sender balance", 500, "Internal server error", nil)
		return
	}

	if err := u.Repository.UpdateUser(receiverUser); err != nil {
		util.Response(c, "Error updating receiver balance", 500, "Internal server error", nil)
		return
	}

	util.Response(c, "Transfer successful", 200, gin.H{
		"sender_balance":    senderUser.InitialBalance,
		"receiver_balance":  receiverUser.InitialBalance,
		"transfer_amount":   transferRequest.Amount,
		"receiver_name":     receiverUser.FirstName + " " + receiverUser.LastName,
		"receiver_email":    receiverUser.Email,
		"receiver_account":  receiverUser.AccountNumber,
	}, nil)
}