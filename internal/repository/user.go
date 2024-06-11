package repository

import "payment-system-four/internal/models"

// FindUserByEmail retrieves a user from the database based on their email address.
func (p *Postgres) FindUserByEmail(email string) (*models.User, error) {
	// Create a new instance of the User model to hold the result of the query
	user := &models.User{}

	// Query the database to find a user with the provided email address
	if err := p.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// If no user is found, return nil for the user and the error
		return nil, err
	}

	// If a user is found, return the user and nil for the error
	return user, nil
}

// CreateUser creates a new user in the database.
func (p *Postgres) CreateUser(user *models.User) error {
	// Create a new user record in the database
	if err := p.DB.Create(user).Error; err != nil {
		// If an error occurs during creation, return the error
		return err
	}

	// Return nil if the creation is successful
	return nil
}

// UpdateUser updates an existing user in the database.
func (p *Postgres) UpdateUser(user *models.User) error {
	// Update the user record in the database
	if err := p.DB.Save(user).Error; err != nil {
		// If an error occurs during update, return the error
		return err
	}

	// Return nil if the update is successful
	return nil
}

func (p *Postgres) TransferFunds(user *models.User, recipient *models.User, amount float64) error {
    tx := p.DB.Begin()
    // deduct the amount from the payer
    user.AvailableBalance -= amount
    // add the amount to the recipient
    recipient.AvailableBalance += amount
    // save the transaction for the payer
    if err := tx.Save(user).Error; err != nil {
        tx.Rollback()
        return err
    }
    // save the transaction for the recipient
    if err := tx.Save(recipient).Error; err != nil {
        tx.Rollback()
        return err
    }
	// save the transaction in the transaction table
    transaction := &models.Transaction{
        PayerAccount:      user.AccountNumber,
        RecipientAccount:  recipient.AccountNumber,
        TransactionType:   "debit",
        TransactionAmount: amount,
        TransactionDate:   time.Now(),
    }
    // save the transaction
    if err := tx.Create(transaction).Error; err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}
func (p *Postgres) FindUserByAccountNumber(account_number int) (*models.User, error) {
    user := &models.User{}
    if err := p.DB.Where("account_no = ?", account_number).First(&user).Error; err != nil {
        return nil, err
    }
    return user, nil
}